/*
 * Compilation:
 *   gcc -Wall -O3 hypervisor.c -o hypervisor
 *
 * */

#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <stdlib.h>

#ifdef __linux__
#include <wait.h>
#else
#include <sys/wait.h>
#endif

#include <errno.h>
#include <time.h>
#include <stdio.h>
#include <stdarg.h>
#include <fcntl.h>
#include <getopt.h>


#define MAXOPT 256
#define INTERVAL 5
#define MAXINTERVAL 180

#define USAGE() do{ \
	fprintf(stderr, "Usage : %s [--daemon] --exec=\"command arg1 arg2 arg3 ...\"\n", argv[0]); \
}while(0)

/* Child process PID, and atomic functions to get and set it.
 * Do not access the internal_child_pid, except using the set_ and get_ functions.
*/
static pid_t   internal_child_pid = 0;
static inline void  set_child_pid(pid_t p) { __atomic_store_n(&internal_child_pid, p, __ATOMIC_SEQ_CST);    }
static inline pid_t get_child_pid(void)    { return __atomic_load_n(&internal_child_pid, __ATOMIC_SEQ_CST); }
static int kill_child_before_exit = 0;

typedef void (*xsa_sigaction)(int, siginfo_t *, void *) ;

static int handle_signal(const int signum, xsa_sigaction sa_handle)
{
    struct sigaction act;

    memset(&act, 0, sizeof act);
    sigemptyset(&act.sa_mask);
    act.sa_sigaction = sa_handle;
    act.sa_flags = SA_SIGINFO | SA_RESTART;

    if (sigaction(signum, &act, NULL))
        return errno;

    return 0;
}

static void handle_sigalrm(int signum, siginfo_t *info, void *context)
{
    if (kill_child_before_exit == 1) {
		const pid_t target = get_child_pid();

		if (target != 0 && info->si_pid != target)
			kill(target, SIGKILL);
	}
    exit(0);
}

static void forward_handler(int signum, siginfo_t *info, void *context)
{
    const pid_t target = get_child_pid();

    if (target != 0 && info->si_pid != target) {
        kill(target, signum);
		kill_child_before_exit = 1;
		handle_signal(SIGALRM, handle_sigalrm);
		// kill child -9 after 2 seconds and exit
		alarm(2);
	} else if(signum == SIGTERM || signum == SIGINT || signum == SIGSTOP ) {
		exit(0);
	}
}

static char *cmd, *cmdopt[MAXOPT + 1];
static int daemonize;

static int parseopt(int argc, char *argv[]);
static int set_nonblock(int fd);
static int getstatus(char *buf, int size, int status);

static int parseopt(int argc, char *argv[])
{
	int ch, i;
	char *token, *tmpptr, *cmdstr;

	cmdstr = cmd = NULL;
	daemonize = 0;
	for(i = 0; i < MAXOPT + 1; i++){
		cmdopt[i] = NULL;
	}

	struct option long_options[] = {
		{"daemon",optional_argument,NULL,'d'},
		{"exec",required_argument,NULL,'e'},
		{0,0,0,0},
	};

	while((ch=getopt_long(argc, argv, "dec:", long_options, NULL)) != -1) {
		switch(ch)
		{
			case 'e':
				if((cmdstr = strdup(optarg)) == NULL )
					return -1;
				break;
			case 'd':
				daemonize = 1;
				break;
			default:
				USAGE();
				return -1;
		}
	}

	if(cmdstr == NULL){
		USAGE();
		return -1;
	}

	for(i = 0;i < MAXOPT + 1;cmdstr = NULL, i++){
		token = strtok_r(cmdstr, " \t\n", &tmpptr);
		if(token == NULL){
			break;
		} else {
			cmdopt[i] = strdup(token);

			if(i == 0){
				cmd = strdup(token);
			}
		}
	}

	if( (cmd == NULL) || (strlen(cmd) == 0) ){
		fprintf(stderr, "Error, cmd should not be empty.\n");
		return -1;
	}

	if(i == MAXOPT + 1){
		fprintf(stderr, "Argument too long\n");
		return -1;
	}

	cmdopt[i] = NULL;

	return 0;
}

static void deal_redirect(char *cmdopt[]) {
	int i;
	char *cmdstr, *openmode;
	FILE *altout=NULL;
	FILE *alterr=NULL;
	
	for(i = 0;i < MAXOPT; i++){
		cmdstr = cmdopt[i];
		if(cmdstr == NULL)
			break;
		
		if(strncmp(cmdstr, "1>", strlen("1>")) == 0) {
			cmdstr +=2;
			openmode = "w";

			if(*cmdstr == '>') {
				cmdstr++;
				openmode = "a";
			}
			if ((altout = fopen(cmdstr, openmode)) != NULL) {
				dup2(fileno(altout), STDOUT_FILENO);
			}
		}

		if(strncmp(cmdstr, "2>", strlen("2>")) == 0) {
			cmdstr +=2;
			openmode = "w";

			if(*cmdstr == '>') {
				cmdstr++;
				openmode = "a";
			} else if(*cmdstr == '&') {
				if(*++cmdstr == '1' && altout != NULL) {
					dup2(fileno(altout), STDERR_FILENO);
					continue;
				}
			}

			if ((alterr = fopen(cmdstr, openmode)) != NULL) {
				dup2(fileno(alterr), STDERR_FILENO);
			}
		}
	}
	if(altout != NULL)
		fclose(altout);
	if(alterr != NULL)
		fclose(alterr);
}

static int set_nonblock(int fd)
{
	int flags = fcntl(fd, F_GETFL, 0);
	if (flags == -1) {
		return -1;
	}
	return fcntl(fd, F_SETFL, flags | O_NONBLOCK);
}

static int getstatus(char *buf, int size, int status)
{
	int n, len;

	len = size;

	if(WIFEXITED(status)){
		n = snprintf(buf, len, "- normal termination, exit status = %d\n", WEXITSTATUS(status));
	} else if(WIFSIGNALED(status)) {
		n = snprintf(buf, len, "- abnormal termination, signal number = %d%s", 
				WTERMSIG(status), 
#ifdef WCOREDUMP
				WCOREDUMP(status) ? "-> core file generated" : "");
#else
		"");
#endif
	} else if(WIFSTOPPED(status)) {
		n = snprintf(buf, len, "child stopped, signal number = %d\n", WSTOPSIG(status));
	}

	return n;
}



void go_daemon() {
	// int fd;

	if (fork() != 0) exit(0); /*  parent exits */
	setsid(); /*  create a new session */

	/*  Every output goes to /dev/null. If Redis is daemonized but
	 *       * the 'logfile' is set to 'stdout' in the configuration file
	 *            * it will not log at all. */
/* 	if ((fd = open("/tmp/mongo4bls.output", O_RDWR, 0)) != -1) {
		dup2(fd, STDIN_FILENO);
		dup2(fd, STDOUT_FILENO);
		dup2(fd, STDERR_FILENO);                                                                                                                        
		if (fd > STDERR_FILENO) close(fd);
	}    
 */
}


int main(int argc, char *argv[])
{
	int ssec = INTERVAL, ret, status;
	int first_start = 1;
	int pipefd[2], waited, alive, isdaemon;
	char buf[1024], info[4096];
	pid_t pid;
	time_t now, last = time(NULL);

	if((ret = parseopt(argc, argv)) < 0 )
		exit(ret);

	daemonize ? go_daemon() : 0;

    /* Install signal forwarders. */
    if (handle_signal(SIGTERM, forward_handler)) {
        fprintf(stderr, "Cannot install signal handlers: %s.\n", strerror(errno));
        return EXIT_FAILURE;
    }

	while(1){
		if(pipe(pipefd) < 0){
			fprintf(stderr, "- make pipe error : %s\n", strerror(errno));
			exit(-1);
		}

		if( (ret = set_nonblock(pipefd[0])) < 0 ){
			fprintf(stderr, "- set read nonblock error : %s\n", strerror(errno));
			exit(-1);
		}

		if((pid = fork()) < 0){
			fprintf(stderr, "- call fork() error : %s\n", strerror(errno));
			exit(-1);
		} else if (pid > 0){
			set_child_pid(pid);
			close(pipefd[1]);
			alive = waited = 1;
			isdaemon = 0;
			while(alive){
				if(waited){
					if(pid != waitpid(pid, &status, 0)){
						sleep(INTERVAL);
						continue;
					} else {
						fprintf(stderr, "- child process[%d] terminated .\n",pid);
						if (first_start && (time(NULL)-last)<=5) {
							fprintf(stderr,"- child process killed in %ld seconds , may wrong ! exit !\n",(time(NULL)-last));
							exit(-1);
						} else
							first_start = 0;
						waited = 0;
					}
				}

				ret = read(pipefd[0], buf, sizeof(buf));
				if(ret < 0){
					if(errno == EAGAIN){
						if(isdaemon == 0){
							fprintf(stderr, "- this daemon process has no output !.\n");
							isdaemon = 1;
						}
						sleep(INTERVAL);
						continue;
					} else {
						fprintf(stderr, "- read pipe error : %s\n", strerror(errno));
						exit(-1);
					}
				} else if(ret == 0) {
					alive = 0;
					close(pipefd[0]);
					fprintf(stderr, "- read zero from pipe of children.\n");
					if(isdaemon == 0){
						getstatus(info, sizeof(info), status);
						fprintf(stderr, "- extra info: %s\n", info);
					} else {
						strcpy(info, "");
					}
					continue;
				} else {
					fprintf(stderr, " - read pipe return: %d bytes\n", ret);
					exit(-1);
				}
			}

			fprintf(stderr, "- process: \"%s\" exit, restart it in %d seconds\n", cmd, ssec);

			sleep(ssec);

			now = time(NULL);
			if(now - last > 3600){
				ssec = INTERVAL;
				last = now;
			} else {
				ssec = (ssec << 1) < MAXINTERVAL ? (ssec << 1) : MAXINTERVAL;
			}
		} else {
			close(pipefd[0]);
			fprintf(stderr, "- execute: \"%s\"\n", cmd);
			deal_redirect(cmdopt);
			if(execvp(cmd, cmdopt) < 0){
				fprintf(stderr, "- execute: \"%s\" error, %s\n", cmd, strerror(errno));
				exit(-1);
			}

		}
	}

	return 0;
}

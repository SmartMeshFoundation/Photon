#ifndef GOICE_HELPER_H
#define GOICE_HELPER_H
#include <stdio.h>
#include <stdlib.h>
#include "pjlib.h"
#include "pjlib-util.h"
#include "pjnath.h"



/* For this demo app, configure longer STUN keep-alive time
* so that it does't clutter the screen output.
*/
#define KA_INTERVAL 300


/* This is our global variables */
 struct app_t
{
	/* Command line options are stored here */
	struct options
	{
		/*
		Component:  A component is a piece of a media stream requiring a
		single transport address; a media stream may require multiple
		components, each of which has to work for the media stream as a
		whole to work.  For media streams based on RTP, there are two
		components per media stream -- one for RTP, and one for RTCP. There may be two component for video transmission
		*/
		unsigned    comp_cnt;
		char*    ns;
		int	    max_host;
		int   regular; //bool
		char*    stun_srv;
		char*    turn_srv;
		int   turn_tcp; //bool
		char*    turn_username;
		char*    turn_password;
		int   turn_fingerprint; //bool
		char *log_file;
		int ka_interval; //interval for stun
	} opt;

	/* Our global variables */
	pj_caching_pool	 cp;
	pj_pool_t		*pool;
	pj_thread_t		*thread;
	pj_bool_t		 thread_quit_flag;
	pj_ice_strans_cfg	 ice_cfg;
	FILE		*log_fhnd;
} ;

typedef  struct IceInstance {
	pj_ice_strans	*icest;
	char name[80];
	void *user_data;
	/* Variables to store parsed remote ICE info */
	struct rem_info
	{
		char		 ufrag[80];
		char		 pwd[80];
		unsigned	 comp_cnt; //ice component number
		pj_sockaddr	 def_addr[PJ_ICE_MAX_COMP]; //remote default address, The first communication address
		unsigned	 cand_cnt; //remote candidate address list
		pj_ice_sess_cand cand[PJ_ICE_ST_MAX_CAND];
	} rem;
	pj_pool_t		* ipool; //pool for instance
}IceInstance;
/* Utility to display error messages */
 void gopjnath_perror(const char *title, pj_status_t status);
 void err_exit(const char *title, pj_status_t status);

#define CHECK(expr)	status=expr; \
			if (status!=PJ_SUCCESS) { \
			    err_exit(#expr, status); \
			}
#define CHECKRETURN(expr)	status=expr; \
			if (status!=PJ_SUCCESS) { \
			    return status; \
			}
#define CHECKBOOLRETURN(expr)	result=expr; \
			if (result!=PJ_TRUE) { \
			    return result; \
			}
/*
* This function checks for events from both timer and ioqueue (for
* network events). It is invoked by the worker thread.
*/
 pj_status_t handle_events(unsigned max_msec, unsigned *p_count);
 int gopjnath_worker_thread(void *unused);
 void cb_on_rx_data(pj_ice_strans *ice_st,
	unsigned comp_id,
	void *pkt, pj_size_t size,
	const pj_sockaddr_t *src_addr,
	unsigned src_addr_len);
 void cb_on_ice_complete(pj_ice_strans *ice_st,
	pj_ice_strans_op op,
	pj_status_t status);
/* log callback to write to file */
 void log_func(int level, const char *data, int len);
 pj_status_t gopjnath_init(void);
IceInstance* gopjnath_create_iceinstance(char*name);
void gopjnath_set_user_data(IceInstance*ii, void *data);
 pj_status_t gopjnath_create_instance(IceInstance*ii, char *name);
 void reset_rem_info(IceInstance*ii);
 void gopjnath_destroy_instance(IceInstance*ii);
 pj_status_t gopjnath_init_session(IceInstance*ii, pj_ice_sess_role role);
 pj_status_t gopjnath_stop_session(IceInstance*ii);

#define PRINT(...)	    \
	printed = pj_ansi_snprintf(p, maxlen - (p-buffer),  \
				   __VA_ARGS__); \
	if (printed <= 0 || printed >= (int)(maxlen - (p-buffer))) \
	    return -PJ_ETOOSMALL; \
	p += printed


/* Utility to create a=candidate SDP attribute */
 int print_cand(char buffer[], unsigned maxlen, const pj_ice_sess_cand *cand);
 int gopjnath_encode_session(IceInstance*ii, char buffer[], int maxlen);
 void gopjnath_show_ice(IceInstance*ii);
 pj_status_t gopjnath_input_remote(IceInstance *ii, char * sdp);
 pj_status_t gopjnath_start_nego(IceInstance*ii);
 pj_status_t gopjnath_send_data(IceInstance*ii, const void *data, pj_size_t len);
 void gopjnath_usage();
pj_status_t goice_init(char* stunsrv, char*turnsrv, char*turnusername, char*turnpassword);
pj_bool_t testmain();
pj_ice_strans	* gopjnath_get_icest(IceInstance*ii);
pj_status_t regThisThread();
#ifdef __linux__ 
//linux code goes here
#define strtok_safe strtok_r
#elif _WIN32
// windows code goes here
#define strtok_safe strtok_s
#else
#define strtok_safe strtok_r
#endif
#endif // !GOICE_HELPER_H


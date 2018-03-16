/* $Id: icedemo.c 4624 2013-10-21 06:37:30Z ming $ */
/*
 * Copyright (C) 2008-2011 Teluu Inc. (http://www.teluu.com)
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 */
#include "helper.h"


#define THIS_FILE   "helper.c"
struct app_t appcfg;
/* Utility to display error messages */
void gopjnath_perror(const char *title, pj_status_t status)
{
	char errmsg[PJ_ERR_MSG_SIZE];

	pj_strerror(status, errmsg, sizeof(errmsg));
	PJ_LOG(1, (THIS_FILE, "%s: %s", title, errmsg));
}

/* Utility: display error message and exit application (usually
 * because of fatal error.
 */
void err_exit(const char *title, pj_status_t status)
{
	if (status != PJ_SUCCESS) {
		gopjnath_perror(title, status);
	}
	PJ_LOG(3, (THIS_FILE, "Shutting down.."));

	//if (appcfg.icest)
	//	pj_ice_strans_destroy(appcfg.icest);

	pj_thread_sleep(500);

	appcfg.thread_quit_flag = PJ_TRUE;
	if (appcfg.thread) {
		pj_thread_join(appcfg.thread);
		pj_thread_destroy(appcfg.thread);
	}

	if (appcfg.ice_cfg.stun_cfg.ioqueue)
		pj_ioqueue_destroy(appcfg.ice_cfg.stun_cfg.ioqueue);

	if (appcfg.ice_cfg.stun_cfg.timer_heap)
		pj_timer_heap_destroy(appcfg.ice_cfg.stun_cfg.timer_heap);

	pj_caching_pool_destroy(&appcfg.cp);

	pj_shutdown();

	if (appcfg.log_fhnd) {
		fclose(appcfg.log_fhnd);
		appcfg.log_fhnd = NULL;
	}

	exit(status != PJ_SUCCESS);
}

/*
 * This function checks for events from both timer and ioqueue (for
 * network events). It is invoked by the worker thread.
 */
pj_status_t handle_events(unsigned max_msec, unsigned *p_count)
{
	enum { MAX_NET_EVENTS = 1 };
	pj_time_val max_timeout = { 0, 0 };
	pj_time_val timeout = { 0, 0 };
	unsigned count = 0, net_event_count = 0;
	int c;

	max_timeout.msec = max_msec;

	/* Poll the timer to run it and also to retrieve the earliest entry. */
	timeout.sec = timeout.msec = 0;
	c = pj_timer_heap_poll(appcfg.ice_cfg.stun_cfg.timer_heap, &timeout);
	if (c > 0)
		count += c;

	/* timer_heap_poll should never ever returns negative value, or otherwise
	 * ioqueue_poll() will block forever!
	 */
	pj_assert(timeout.sec >= 0 && timeout.msec >= 0);
	if (timeout.msec >= 1000) timeout.msec = 999;

	/* compare the value with the timeout to wait from timer, and use the
	 * minimum value.
	*/
	if (PJ_TIME_VAL_GT(timeout, max_timeout))
		timeout = max_timeout;

	/* Poll ioqueue.
	 * Repeat polling the ioqueue while we have immediate events, because
	 * timer heap may process more than one events, so if we only process
	 * one network events at a time (such as when IOCP backend is used),
	 * the ioqueue may have trouble keeping up with the request rate.
	 *
	 * For example, for each send() request, one network event will be
	 *   reported by ioqueue for the send() completion. If we don't poll
	 *   the ioqueue often enough, the send() completion will not be
	 *   reported in timely manner.
	 */
	do {
		c = pj_ioqueue_poll(appcfg.ice_cfg.stun_cfg.ioqueue, &timeout);
		if (c < 0) {
			pj_status_t err = pj_get_netos_error();
			pj_thread_sleep(PJ_TIME_VAL_MSEC(timeout));
			if (p_count)
				*p_count = count;
			return err;
		}
		else if (c == 0) {
			break;
		}
		else {
			net_event_count += c;
			timeout.sec = timeout.msec = 0;
		}
	} while (c > 0 && net_event_count < MAX_NET_EVENTS);

	count += net_event_count;
	if (p_count)
		*p_count = count;

	return PJ_SUCCESS;

}

/*
 * This is the worker thread that polls event in the background.
 */
int gopjnath_worker_thread(void *unused)
{
	PJ_UNUSED_ARG(unused);

	while (!appcfg.thread_quit_flag) {
		handle_events(500, NULL);
	}

	return 0;
}
extern void  ice_cb(pj_ice_strans *ice_strans, pj_ice_strans_op op, pj_status_t status);
extern void  data_cb(pj_ice_strans *ice_st, unsigned comp_id, void *pkt, pj_size_t size, const pj_sockaddr_t *src_addr, unsigned src_addr_len);

/* log callback to write to file */
void log_func(int level, const char *data, int len)
{
	pj_log_write(level, data, len);
	if (appcfg.log_fhnd) {
		if (fwrite(data, len, 1, appcfg.log_fhnd) != 1)
			return;
	}
}

/*
 * This is the main application initialization function. It is called
 * once (and only once) during application initialization sequence by
 * main().
 */
pj_status_t gopjnath_init(void)
{
	pj_status_t status;
	pj_time_val now;
	pj_str_t tmpstr;
	if (appcfg.opt.log_file) {
		appcfg.log_fhnd = fopen(appcfg.opt.log_file, "a");
		pj_log_set_log_func(&log_func);
	}


	/* Initialize the libraries before anything else */
	CHECKRETURN(pj_init());
	CHECKRETURN(pj_gettimeofday(&now));
	pj_srand(now.sec);
	CHECKRETURN(pjlib_util_init());
	CHECKRETURN(pjnath_init());

	/* Must create pool factory, where memory allocations come from */
	pj_caching_pool_init(&appcfg.cp, NULL, 0);

	/* Init our ICE settings with null values */
	pj_ice_strans_cfg_default(&appcfg.ice_cfg);

	appcfg.ice_cfg.stun_cfg.pf = &appcfg.cp.factory;

	/* Create application memory pool */
	appcfg.pool = pj_pool_create(&appcfg.cp.factory, "gopjnath",
		512, 512, NULL);

	/* Create timer heap for timer stuff */
	CHECKRETURN(pj_timer_heap_create(appcfg.pool, 100,
		&appcfg.ice_cfg.stun_cfg.timer_heap));

	/* and create ioqueue for network I/O stuff */
	CHECKRETURN(pj_ioqueue_create(appcfg.pool, 16,
		&appcfg.ice_cfg.stun_cfg.ioqueue));

	/* something must poll the timer heap and ioqueue,
	 * unless we're on Symbian where the timer heap and ioqueue run
	 * on themselves.
	 */
	CHECKRETURN(pj_thread_create(appcfg.pool, "gopjnath", &gopjnath_worker_thread, NULL, 0, 0, &appcfg.thread));

	appcfg.ice_cfg.af = pj_AF_INET();

	/* Create DNS resolver if nameserver is set */
	if (appcfg.opt.ns) {
		CHECKRETURN(pj_dns_resolver_create(&appcfg.cp.factory,
			"resolver",
			0,
			appcfg.ice_cfg.stun_cfg.timer_heap,
			appcfg.ice_cfg.stun_cfg.ioqueue,
			&appcfg.ice_cfg.resolver));
		tmpstr = pj_str(appcfg.opt.ns);
		CHECKRETURN(pj_dns_resolver_set_ns(appcfg.ice_cfg.resolver, 1,
			&tmpstr, NULL));
	}

	/* -= Start initializing ICE stream transport config =- */

	/* Maximum number of host candidates */
	if (appcfg.opt.max_host != -1)
		appcfg.ice_cfg.stun.max_host_cands = appcfg.opt.max_host;

	/* Nomination strategy */
	if (appcfg.opt.regular)
		appcfg.ice_cfg.opt.aggressive = PJ_FALSE;
	else
		appcfg.ice_cfg.opt.aggressive = PJ_TRUE;

	/* Configure STUN/srflx candidate resolution */
	if (appcfg.opt.stun_srv) {
		char *pos;
		tmpstr = pj_str(appcfg.opt.stun_srv);
		/* Command line option may contain port number */
		if ((pos = pj_strchr(&tmpstr, ':')) != NULL) {
			appcfg.ice_cfg.stun.server.ptr = appcfg.opt.stun_srv;
			appcfg.ice_cfg.stun.server.slen = (pos - tmpstr.ptr);

			appcfg.ice_cfg.stun.port = (pj_uint16_t)atoi(pos + 1);
		}
		else {
			appcfg.ice_cfg.stun.server = tmpstr;
			appcfg.ice_cfg.stun.port = PJ_STUN_PORT;
		}

		/* For this demo app, configure longer STUN keep-alive time
		 * so that it does't clutter the screen output.
		 */
		appcfg.ice_cfg.stun.cfg.ka_interval = appcfg.opt.ka_interval;
	}

	/* Configure TURN candidate */
	if (appcfg.opt.turn_srv) {
		char *pos;
		tmpstr = pj_str(appcfg.opt.turn_srv);
		/* Command line option may contain port number */
		if ((pos = pj_strchr(&tmpstr, ':')) != NULL) {
			appcfg.ice_cfg.turn.server.ptr = tmpstr.ptr;
			appcfg.ice_cfg.turn.server.slen = (pos - tmpstr.ptr);

			appcfg.ice_cfg.turn.port = (pj_uint16_t)atoi(pos + 1);
		}
		else {
			appcfg.ice_cfg.turn.server = tmpstr;
			appcfg.ice_cfg.turn.port = PJ_STUN_PORT;
		}

		/* TURN credential */
		appcfg.ice_cfg.turn.auth_cred.type = PJ_STUN_AUTH_CRED_STATIC;
		tmpstr = pj_str(appcfg.opt.turn_username);
		appcfg.ice_cfg.turn.auth_cred.data.static_cred.username = tmpstr;
		appcfg.ice_cfg.turn.auth_cred.data.static_cred.data_type = PJ_STUN_PASSWD_PLAIN;
		tmpstr = pj_str(appcfg.opt.turn_password);
		appcfg.ice_cfg.turn.auth_cred.data.static_cred.data = tmpstr;

		/* Connection type to TURN server */
		if (appcfg.opt.turn_tcp)
			appcfg.ice_cfg.turn.conn_type = PJ_TURN_TP_TCP;
		else
			appcfg.ice_cfg.turn.conn_type = PJ_TURN_TP_UDP;

		/* For this demo app, configure longer keep-alive time
		 * so that it does't clutter the screen output.
		 */
		appcfg.ice_cfg.turn.alloc_param.ka_interval = appcfg.opt.ka_interval;
	}

	/* -= That's it for now, initialization is complete =- */
	return PJ_SUCCESS;
}
IceInstance* gopjnath_create_iceinstance(char*name) {
	IceInstance*ii = (IceInstance*)calloc(sizeof(IceInstance), 1);
	/* Create application memory pool */
	ii->ipool = pj_pool_create(&appcfg.cp.factory, "instancepool", 512, 512, NULL);
	strncpy(ii->name, name, 80);
	return ii;
}
void gopjnath_set_user_data(IceInstance*ii, void *data) {
	ii->user_data = data;
}
/*
 * Create ICE stream transport instance, invoked from the menu.
 */
pj_status_t gopjnath_create_instance(IceInstance*ii, char *name)
{
	pj_ice_strans_cb icecb;
	pj_status_t status;


	/* init the callback */
	pj_bzero(&icecb, sizeof(icecb));
	icecb.on_rx_data = data_cb;
	icecb.on_ice_complete = ice_cb;

	/* create the instance */
	status = pj_ice_strans_create(name,		    /* object name  */
		&appcfg.ice_cfg,	    /* settings	    */
		appcfg.opt.comp_cnt,	    /* comp_cnt	    */
		ii->user_data,			    /* user data    */
		&icecb,			    /* callback	    */
		&ii->icest)		    /* instance ptr */
		;
	if (status != PJ_SUCCESS)
	{
		gopjnath_perror("error creating ice", status);
		return status;
	}
	else {
		PJ_LOG(3, (THIS_FILE, "%s ICE instance successfully created", ii->name));
		return PJ_SUCCESS;
	}
}

/* Utility to nullify parsed remote info */
void reset_rem_info(IceInstance*ii)
{
	pj_bzero(&ii->rem, sizeof(ii->rem));
}


/*
 * Destroy ICE stream transport instance, invoked from the menu.
 */
void gopjnath_destroy_instance(IceInstance*ii)
{
	if (ii->icest != NULL) {
		gopjnath_stop_session(ii);
		pj_ice_strans_destroy(ii->icest);
		ii->icest = NULL;
	}
	else {
		PJ_LOG(3, (THIS_FILE, "%s Error: No ICE instance, create it first", ii->name));
	}
	if (ii->user_data != NULL) {
		free(ii->user_data);
	}
	if (ii->ipool != NULL) {
		pj_pool_release(ii->ipool);
	}
	PJ_LOG(3, (THIS_FILE, "%s ICE instance destroyed", ii->name));
	free(ii);
	return;
}


/*
 * Create ICE session, invoked from the menu.
 */
pj_status_t gopjnath_init_session(IceInstance*ii, pj_ice_sess_role role)
{
	pj_status_t status;
	PJ_LOG(3, (THIS_FILE, "%s gopjnath_init_session", ii->name));
	if (ii->icest == NULL) {
		PJ_LOG(1, (THIS_FILE, "%s  gopjnath_init_session Error: No ICE instance, create it first", ii->name));
		return PJ_EINVAL;
	}

	if (pj_ice_strans_has_sess(ii->icest)) {
		PJ_LOG(1, (THIS_FILE, "%s gopjnath_init_session Error: Session already created", ii->name));
		return PJ_EINVAL;
	}

	status = pj_ice_strans_init_ice(ii->icest, role, NULL, NULL);
	if (status != PJ_SUCCESS) {
		gopjnath_perror("error creating session", status);
		return status;
	}

	else {
		PJ_LOG(3, (THIS_FILE, "%s ICE session created", ii->name));
	}
	reset_rem_info(ii);
	return PJ_SUCCESS;
}


/*
 * Stop/destroy ICE session, invoked from the menu.
 */
pj_status_t gopjnath_stop_session(IceInstance*ii)
{
	pj_status_t status;

	if (ii->icest == NULL) {
		PJ_LOG(1, (THIS_FILE, "%s gopjnath_stop_session Error: No ICE instance, create it first", ii->name));
		return PJ_EINVAL;
	}

	if (!pj_ice_strans_has_sess(ii->icest)) {
		PJ_LOG(1, (THIS_FILE, "%s gopjnath_stop_session Error: No ICE session, initialize first", ii->name));
		return PJ_EINVAL;
	}

	status = pj_ice_strans_stop_ice(ii->icest);
	if (status != PJ_SUCCESS) {
		gopjnath_perror("error stopping session", status);
		return status;
	}
	else
		PJ_LOG(3, (THIS_FILE, "%s ICE session stopped", ii->name));

	reset_rem_info(ii);
	return PJ_SUCCESS;
}

/* Utility to create a=candidate SDP attribute */
int print_cand(char buffer[], unsigned maxlen,
	const pj_ice_sess_cand *cand)
{
	char ipaddr[PJ_INET6_ADDRSTRLEN];
	char *p = buffer;
	int printed;

	PRINT("a=candidate:%.*s %u UDP %u %s %u typ ",
		(int)cand->foundation.slen,
		cand->foundation.ptr,
		(unsigned)cand->comp_id,
		cand->prio,
		pj_sockaddr_print(&cand->addr, ipaddr,
			sizeof(ipaddr), 0),
			(unsigned)pj_sockaddr_get_port(&cand->addr));

	PRINT("%s\n",
		pj_ice_get_cand_type_name(cand->type));

	if (p == buffer + maxlen)
		return -PJ_ETOOSMALL;

	*p = '\0';

	return (int)(p - buffer);
}

/*
 * Encode ICE information in SDP.
 */
int gopjnath_encode_session(IceInstance*ii, char buffer[], int maxlen)
{
	char *p = buffer;
	unsigned comp;
	int printed;
	pj_str_t local_ufrag, local_pwd;
	pj_status_t status;

	/* Write "dummy" SDP v=, o=, s=, and t= lines */
	PRINT("v=0\no=- 3414953978 3414953978 IN IP4 localhost\ns=ice\nt=0 0\n");

	/* Get ufrag and pwd from current session */
	pj_ice_strans_get_ufrag_pwd(ii->icest, &local_ufrag, &local_pwd,
		NULL, NULL);

	/* Write the a=ice-ufrag and a=ice-pwd attributes */
	PRINT("a=ice-ufrag:%.*s\na=ice-pwd:%.*s\n",
		(int)local_ufrag.slen,
		local_ufrag.ptr,
		(int)local_pwd.slen,
		local_pwd.ptr);

	/* Write each component */
	for (comp = 0; comp < appcfg.opt.comp_cnt; ++comp) {
		unsigned j, cand_cnt;
		pj_ice_sess_cand cand[PJ_ICE_ST_MAX_CAND];
		char ipaddr[PJ_INET6_ADDRSTRLEN];

		/* Get default candidate for the component */
		status = pj_ice_strans_get_def_cand(ii->icest, comp + 1, &cand[0]);
		if (status != PJ_SUCCESS)
			return -status;

		/* Write the default address */
		if (comp == 0) {
			/* For component 1, default address is in m= and c= lines */
			PRINT("m=audio %d RTP/AVP 0\n"
				"c=IN IP4 %s\n",
				(int)pj_sockaddr_get_port(&cand[0].addr),
				pj_sockaddr_print(&cand[0].addr, ipaddr,
					sizeof(ipaddr), 0));
		}
		else if (comp == 1) {
			/* For component 2, default address is in a=rtcp line */
			PRINT("a=rtcp:%d IN IP4 %s\n",
				(int)pj_sockaddr_get_port(&cand[0].addr),
				pj_sockaddr_print(&cand[0].addr, ipaddr,
					sizeof(ipaddr), 0));
		}
		else {
			/* For other components, we'll just invent this.. */
			PRINT("a=Xice-defcand:%d IN IP4 %s\n",
				(int)pj_sockaddr_get_port(&cand[0].addr),
				pj_sockaddr_print(&cand[0].addr, ipaddr,
					sizeof(ipaddr), 0));
		}

		/* Enumerate all candidates for this component */
		cand_cnt = PJ_ARRAY_SIZE(cand);
		status = pj_ice_strans_enum_cands(ii->icest, comp + 1,
			&cand_cnt, cand);
		if (status != PJ_SUCCESS)
			return -status;

		/* And encode the candidates as SDP */
		for (j = 0; j < cand_cnt; ++j) {
			printed = print_cand(p, maxlen - (unsigned)(p - buffer), &cand[j]);
			if (printed < 0)
				return -PJ_ETOOSMALL;
			p += printed;
		}
	}

	if (p == buffer + maxlen)
		return -PJ_ETOOSMALL;

	*p = '\0';
	return (int)(p - buffer);
}


/*
 * Show information contained in the ICE stream transport. This is
 * invoked from the menu.
 */
void gopjnath_show_ice(IceInstance*ii)
{
	static char buffer[1000];
	int len;

	if (ii->icest == NULL) {
		PJ_LOG(1, (THIS_FILE, "%s gopjnath_show_ice Error: No ICE instance, create it first", ii->name));
		return;
	}
	puts(ii->name);
	puts("General info");
	puts("---------------");
	printf("Component count    : %d\n", appcfg.opt.comp_cnt);
	printf("Status             : ");
	if (pj_ice_strans_sess_is_complete(ii->icest))
		puts("negotiation complete");
	else if (pj_ice_strans_sess_is_running(ii->icest))
		puts("negotiation is in progress");
	else if (pj_ice_strans_has_sess(ii->icest))
		puts("session ready");
	else
		puts("session not created");

	if (!pj_ice_strans_has_sess(ii->icest)) {
		puts("Create the session first to see more info");
		return;
	}

	printf("Negotiated comp_cnt: %d\n",
		pj_ice_strans_get_running_comp_cnt(ii->icest));
	printf("Role               : %s\n",
		pj_ice_strans_get_role(ii->icest) == PJ_ICE_SESS_ROLE_CONTROLLED ?
		"controlled" : "controlling");

	len = gopjnath_encode_session(ii, buffer, sizeof(buffer));
	if (len < 0)
		err_exit("not enough buffer to show ICE status", -len);

	puts("");
	printf("Local SDP (paste this to remote host):\n"
		"--------------------------------------\n"
		"%s\n", buffer);


	puts("");
	puts("Remote info:\n"
		"----------------------");
	if (ii->rem.cand_cnt == 0) {
		puts("No remote info yet");
	}
	else {
		unsigned i;

		printf("Remote ufrag       : %s\n", ii->rem.ufrag);
		printf("Remote password    : %s\n", ii->rem.pwd);
		printf("Remote cand. cnt.  : %d\n", ii->rem.cand_cnt);

		for (i = 0; i < ii->rem.cand_cnt; ++i) {
			len = print_cand(buffer, sizeof(buffer), &ii->rem.cand[i]);
			if (len < 0)
				err_exit("not enough buffer to show ICE status", -len);

			printf("  %s", buffer);
		}
	}
}


/*
 * Input and parse SDP from the remote (containing remote's ICE information)
 * and save it to global variables.
 */
pj_status_t gopjnath_input_remote(IceInstance *ii, char * sdp)
{
	unsigned media_cnt = 0;
	unsigned comp0_port = 0;
	char     comp0_addr[80];
	char * line = NULL, *savep = NULL;

	reset_rem_info(ii);

	comp0_addr[0] = '\0';

	while (1) {
		if (savep == NULL) {
			line = strtok_safe(sdp, "\n", &savep);
		}
		else {
			line = strtok_safe(NULL, "\n", &savep);
		}
		if (line == NULL) {
			break;
		}
		PJ_LOG(4, (THIS_FILE, "%s processing line :%s", ii->name, line));
		pj_size_t len;

		len = strlen(line);
		while (len && (line[len - 1] == '\r' || line[len - 1] == '\n'))
			line[--len] = '\0';

		while (len && pj_isspace(*line))
			++line, --len;

		if (len == 0)
			break;

		/* Ignore subsequent media descriptors */
		if (media_cnt > 1)
			continue;

		switch (line[0]) {
		case 'm':
		{
			int cnt;
			char media[32], portstr[32];

			++media_cnt;
			if (media_cnt > 1) {
				puts("Media line ignored");
				break;
			}

			cnt = sscanf(line + 2, "%s %s RTP/", media, portstr);
			if (cnt != 2) {
				PJ_LOG(1, (THIS_FILE, "%s Error parsing media line", ii->name));
				goto on_error;
			}

			comp0_port = atoi(portstr);

		}
		break;
		case 'c':
		{
			int cnt;
			char c[32], net[32], ip[80];

			cnt = sscanf(line + 2, "%s %s %s", c, net, ip);
			if (cnt != 3) {
				PJ_LOG(1, (THIS_FILE, "%s Error parsing connection line", ii->name));
				goto on_error;
			}

			strcpy(comp0_addr, ip);
		}
		break;
		case 'a':
		{
			char *attr = strtok(line + 2, ": \t\r\n");
			if (strcmp(attr, "ice-ufrag") == 0) {
				strcpy(ii->rem.ufrag, attr + strlen(attr) + 1);
			}
			else if (strcmp(attr, "ice-pwd") == 0) {
				strcpy(ii->rem.pwd, attr + strlen(attr) + 1);
			}
			else if (strcmp(attr, "rtcp") == 0) {
				char *val = attr + strlen(attr) + 1;
				int af, cnt;
				int port;
				char net[32], ip[64];
				pj_str_t tmp_addr;
				pj_status_t status;

				cnt = sscanf(val, "%d IN %s %s", &port, net, ip);
				if (cnt != 3) {
					PJ_LOG(1, (THIS_FILE, "%s Error parsing rtcp attribute", ii->name));
					goto on_error;
				}

				if (strchr(ip, ':'))
					af = pj_AF_INET6();
				else
					af = pj_AF_INET();

				pj_sockaddr_init(af, &ii->rem.def_addr[1], NULL, 0);
				tmp_addr = pj_str(ip);
				status = pj_sockaddr_set_str_addr(af, &ii->rem.def_addr[1],
					&tmp_addr);
				if (status != PJ_SUCCESS) {
					PJ_LOG(1, (THIS_FILE, "%s Invalid IP address", ii->name));
					goto on_error;
				}
				pj_sockaddr_set_port(&ii->rem.def_addr[1], (pj_uint16_t)port);

			}
			else if (strcmp(attr, "candidate") == 0) {
				char *sdpcand = attr + strlen(attr) + 1;
				int af, cnt;
				char foundation[32], transport[12], ipaddr[80], type[32];
				pj_str_t tmpaddr;
				int comp_id, prio, port;
				pj_ice_sess_cand *cand;
				pj_status_t status;
				cnt = sscanf(sdpcand, "%s %d %s %d %s %d typ %s",
					foundation,
					&comp_id,
					transport,
					&prio,
					ipaddr,
					&port,
					type);
				if (cnt != 7) {
					PJ_LOG(1, (THIS_FILE, "%s error: Invalid ICE candidate line", ii->name));
					goto on_error;
				}
				cand = &ii->rem.cand[ii->rem.cand_cnt];
				pj_bzero(cand, sizeof(*cand));

				if (strcmp(type, "host") == 0)
					cand->type = PJ_ICE_CAND_TYPE_HOST;
				else if (strcmp(type, "srflx") == 0)
					cand->type = PJ_ICE_CAND_TYPE_SRFLX;
				else if (strcmp(type, "relay") == 0)
					cand->type = PJ_ICE_CAND_TYPE_RELAYED;
				else {
					PJ_LOG(1, (THIS_FILE, "%s Error: invalid candidate type '%s'", ii->name,
						type));
					goto on_error;
				}
				cand->comp_id = (pj_uint8_t)comp_id;
				pj_strdup2(ii->ipool, &cand->foundation, foundation);
				cand->prio = prio;

				if (strchr(ipaddr, ':'))
					af = pj_AF_INET6();
				else
					af = pj_AF_INET();
				tmpaddr = pj_str(ipaddr);
				pj_sockaddr_init(af, &cand->addr, NULL, 0);
				status = pj_sockaddr_set_str_addr(af, &cand->addr, &tmpaddr);
				if (status != PJ_SUCCESS) {
					PJ_LOG(1, (THIS_FILE, "%s Error: invalid IP address '%s'", ii->name,
						ipaddr));
					goto on_error;
				}
				pj_sockaddr_set_port(&cand->addr, (pj_uint16_t)port);

				++ii->rem.cand_cnt;
				if (cand->comp_id > ii->rem.comp_cnt)
					ii->rem.comp_cnt = cand->comp_id;
			}
		}
		break;
		}
	}

	if (ii->rem.cand_cnt == 0 ||
		ii->rem.ufrag[0] == 0 ||
		ii->rem.pwd[0] == 0 ||
		ii->rem.comp_cnt == 0)
	{
		PJ_LOG(1, (THIS_FILE, "%s Error: not enough info", ii->name));
		goto on_error;
	}

	if (comp0_port == 0 || comp0_addr[0] == '\0') {
		PJ_LOG(1, (THIS_FILE, "%s Error: default address for component 0 not found", ii->name));
		goto on_error;
	}
	else {
		int af;
		pj_str_t tmp_addr;
		pj_status_t status;

		if (strchr(comp0_addr, ':'))
			af = pj_AF_INET6();
		else
			af = pj_AF_INET();

		pj_sockaddr_init(af, &ii->rem.def_addr[0], NULL, 0);
		tmp_addr = pj_str(comp0_addr);
		status = pj_sockaddr_set_str_addr(af, &ii->rem.def_addr[0],
			&tmp_addr);
		if (status != PJ_SUCCESS) {
			PJ_LOG(1, (THIS_FILE, "%s Invalid IP address in c= line", ii->name));
			goto on_error;
		}
		pj_sockaddr_set_port(&ii->rem.def_addr[0], (pj_uint16_t)comp0_port);
	}

	PJ_LOG(3, (THIS_FILE, "%s Done, %d remote candidate(s) added", ii->name,
		ii->rem.cand_cnt));
	return PJ_SUCCESS;

on_error:
	reset_rem_info(ii);
	return PJ_EINVAL;
}

/*
 * Start ICE negotiation! This function is invoked from the menu.
 */
pj_status_t gopjnath_start_nego(IceInstance*ii)
{
	pj_str_t rufrag, rpwd;
	pj_status_t status;

	if (ii->icest == NULL) {
		PJ_LOG(1, (THIS_FILE, "%s gopjnath_start_nego Error: No ICE instance, create it first", ii->name));
		return PJ_EINVAL;
	}

	if (!pj_ice_strans_has_sess(ii->icest)) {
		PJ_LOG(1, (THIS_FILE, "%s gopjnath_start_nego Error: No ICE session, initialize first", ii->name));
		return  PJ_EINVAL;
	}

	if (ii->rem.cand_cnt == 0) {
		PJ_LOG(1, (THIS_FILE, "%s gopjnath_start_nego Error: No remote info, input remote info first", ii->name));
		return  PJ_EINVAL;
	}

	PJ_LOG(3, (THIS_FILE, "%s gopjnath_start_nego Starting ICE negotiation..", ii->name));

	status = pj_ice_strans_start_ice(ii->icest,
		pj_cstr(&rufrag, ii->rem.ufrag),
		pj_cstr(&rpwd, ii->rem.pwd),
		ii->rem.cand_cnt,
		ii->rem.cand);
	if (status != PJ_SUCCESS) {
		gopjnath_perror("Error starting ICE", status);
		return status;
	}
	else {
		PJ_LOG(3, (THIS_FILE, "%s ICE negotiation started", ii->name));
		return PJ_SUCCESS;
	}
}


/*
 * Send application data to remote agent.
 */
pj_status_t gopjnath_send_data(IceInstance*ii, const void *data, pj_size_t len)
{
	pj_status_t status;
	unsigned int comp_id = 1;

	if (ii->icest == NULL) {
		PJ_LOG(1, (THIS_FILE, "%s Error: No ICE instance, create it first", ii->name));
		return PJ_EINVAL;
	}

	if (!pj_ice_strans_has_sess(ii->icest)) {
		PJ_LOG(1, (THIS_FILE, "%s Error: No ICE session, initialize first", ii->name));
		return PJ_EINVAL;
	}


	if (!pj_ice_strans_sess_is_complete(ii->icest)) {
		PJ_LOG(1, (THIS_FILE, "%s Error: ICE negotiation has not been started or is in progress", ii->name));
		return PJ_EINVAL;
	}
	status = regThisThread();
	if (status != PJ_SUCCESS) {
		PJ_LOG(1, (THIS_FILE, "%s gopjnath_send_data reg this thread err ", ii->name));
		return status;
	}
	/*
		if (comp_id<1 || comp_id>pj_ice_strans_get_running_comp_cnt(appcfg.icest)) {
			PJ_LOG(1, (THIS_FILE, "Error: invalid component ID"));
			return;
		}*/
	PJ_LOG(3, (THIS_FILE, "%s Send data %d", ii->name, len));
	status = pj_ice_strans_sendto(ii->icest, comp_id, data, len,
		&ii->rem.def_addr[comp_id - 1],
		pj_sockaddr_get_len(&ii->rem.def_addr[comp_id - 1]));
	if (status != PJ_SUCCESS)
	{
		gopjnath_perror("Error sending data", status);
		return status;
	}
	else {
		PJ_LOG(3, (THIS_FILE, "%s Data sent", ii->name));
		return PJ_SUCCESS;
	}
}


pj_status_t goice_init(char* stunsrv, char*turnsrv, char*turnusername, char*turnpassword) {
	pj_status_t status;
	appcfg.opt.comp_cnt = 1;
	appcfg.opt.max_host = -1;
	appcfg.opt.ka_interval = 300;
	appcfg.opt.log_file = NULL;
	appcfg.opt.ns = NULL;
	appcfg.opt.regular = 1;
	appcfg.opt.stun_srv = stunsrv;
	appcfg.opt.turn_fingerprint = 0;
	appcfg.opt.turn_srv = turnsrv;
	appcfg.opt.turn_username = turnusername;
	appcfg.opt.turn_password = turnpassword;
	appcfg.opt.turn_tcp = 0;

	status = gopjnath_init();
	return status;
}
void teststrtok() {
	char *str = "abc.def.ghi.mmm";
	pj_ssize_t found_idx;
	pj_str_t in_str = pj_str(str);
	pj_str_t token, delim;

	while (*str && !pj_isdigit(*str))
		str++;

	delim = pj_str(".-");
	for (found_idx = pj_strtok(&in_str, &delim, &token, 0);
		found_idx != in_str.slen;
		found_idx = pj_strtok(&in_str, &delim, &token,
			found_idx + token.slen))
	{
		printf("token=%s,found_idx=%d\n", token.ptr, found_idx);
		printf("in_str=%s\n", in_str.ptr);
	}

}
pj_bool_t testmain() {
	pj_status_t status;
	pj_bool_t result;
	appcfg.opt.comp_cnt = 1;
	appcfg.opt.max_host = -1;
	appcfg.opt.ka_interval = 300;
	appcfg.opt.log_file = NULL;
	appcfg.opt.ns = NULL;
	appcfg.opt.regular = 1;
	appcfg.opt.stun_srv = "182.254.155.208:3478";
	appcfg.opt.turn_fingerprint = 0;
	appcfg.opt.turn_srv = "182.254.155.208:3478";
	appcfg.opt.turn_username = "bai";
	appcfg.opt.turn_password = "bai";
	appcfg.opt.turn_tcp = 0;
	teststrtok();
	return PJ_TRUE;
	status = gopjnath_init();
	if (status != PJ_SUCCESS)
		return PJ_FALSE;
	IceInstance iis;
	pj_bzero(&iis, sizeof(iis));
	CHECKBOOLRETURN(gopjnath_create_instance(&iis, "test"));
	//CHECKBOOLRETURN(gopjnath_init_session(&iis, 'o'));
	pj_thread_sleep(1000);;
	gopjnath_input_remote(&iis, "");
	gopjnath_show_ice(&iis);
	CHECKBOOLRETURN(gopjnath_stop_session(&iis));
	gopjnath_destroy_instance(&iis);

	//err_exit("Quitting..", PJ_SUCCESS);
	return PJ_TRUE;
}
pj_status_t regThisThread() {
	pj_thread_desc rtpdesc;
	pj_thread_t *thread = NULL;
	if (!pj_thread_is_registered())
	{
		return pj_thread_register(NULL, rtpdesc, &thread);
	}
	return PJ_SUCCESS;
}

pj_ice_strans	* gopjnath_get_icest(IceInstance*ii) {
	return ii->icest;
}

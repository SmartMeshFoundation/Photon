package xmpp

import (
	"fmt"
)

func (c *Client) PingC2S(jid, server string) error {
	if jid == "" {
		jid = c.jid
	}
	if server == "" {
		server = c.domain
	}
	_, err := fmt.Fprintf(c.w, "<iq from='%s' to='%s' id='c2s1' type='get'>\n"+
		"<ping xmlns='urn:xmpp:ping'/>\n"+
		"</iq>",
		xmlEscape(jid), xmlEscape(server))
	return err
}

func (c *Client) PingS2S(fromServer, toServer string) error {
	_, err := fmt.Fprintf(c.w, "<iq from='%s' to='%s' id='s2s1' type='get'>\n"+
		"<ping xmlns='urn:xmpp:ping'/>\n"+
		"</iq>",
		xmlEscape(fromServer), xmlEscape(toServer))
	return err
}

func (c *Client) SendResultPing(id, toServer string) error {
	_, err := fmt.Fprintf(c.w, "<iq type='result' to='%s' id='%s'/>",
		xmlEscape(toServer), xmlEscape(id))
	return err
}

func (c*Client) SendOnlinePing(id,from,to string) error{
	_, err := fmt.Fprintf(c.w, "<iq from='%s' to='%s' id='%s' type='get'>\n"+
		"<ping xmlns='urn:xmpp:ping'/>\n"+
		"</iq>",
		xmlEscape(from), xmlEscape(to),xmlEscape(id))
	return err
}
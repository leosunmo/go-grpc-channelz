package channelz

import (
	"context"
	"fmt"
	"io"

	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	log "google.golang.org/grpc/grpclog"
)

func (h *grpcChannelzHandler) WriteSocketPage(w io.Writer, socket int64) {
	writeHeader(w, fmt.Sprintf("ChannelZ socket %d", socket))
	h.writeSocket(w, socket)
	writeFooter(w)
}

// writeSocket writes HTML to w containing RPC single socket stats.
//
// It includes neither a header nor footer, so you can embed this data in other pages.
func (h *grpcChannelzHandler) writeSocket(w io.Writer, socket int64) {
	if err := socketTemplate.Execute(w, h.getSocket(socket)); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func (h *grpcChannelzHandler) getSocket(socketID int64) *channelzgrpc.GetSocketResponse {
	client, err := h.connect()
	if err != nil {
		log.Errorf("Error creating channelz client %+v", err)
		return nil
	}
	ctx := context.Background()
	socket, err := client.GetSocket(ctx, &channelzgrpc.GetSocketRequest{SocketId: socketID})
	if err != nil {
		log.Errorf("Error querying GetSocket %+v", err)
		return nil
	}
	return socket
}

const socketTemplateHTML = `
<table frame=box cellspacing=0 cellpadding=2 class="vertical">
	<tr>
		<th>SocketId</th>
		<td>
			{{.Socket.Ref.SocketId}}
		</td>
	</tr>
	<tr>
		<th>Socket Name</th>
		<td>
			{{.Socket.Ref.Name}}
		</td>
	</tr>
	<tr>
		<th>Socket Local -> Remote</th>
		<td>
			<pre>{{ipToString .Socket.Local}} -> {{ipToString .Socket.Remote}} {{with .Socket.RemoteName}}({{.}}){{end}}</pre>
		</td>
	</tr>
	<tr>
		<th>StreamsStarted</th>
		<td>{{.Socket.Data.StreamsStarted}}</td>
	</tr>
	<tr>
		<th>StreamsSucceeded</th>
		<td>{{.Socket.Data.StreamsSucceeded}}</td>
	</tr>
	<tr>
		<th>StreamsFailed</th>
		<td>{{.Socket.Data.StreamsFailed}}</td>
	</tr>
	<tr>
		<th>MessagesSent</th>
		<td>{{.Socket.Data.MessagesSent}}</td>
	</tr>
	<tr>
		<th>MessagesReceived</th>
		<td>{{.Socket.Data.MessagesReceived}}</td>
	</tr>
	<tr>
		<th>KeepAlivesSent</th>
		<td>{{.Socket.Data.KeepAlivesSent}}</td>
	</tr>
	<tr>
		<th>LastLocalStreamCreated</th>
		<td>{{.Socket.Data.LastLocalStreamCreatedTimestamp | timestamp}}</td>
	</tr>
	<tr>
		<th>LastRemoteStreamCreated</th>
		<td>{{.Socket.Data.LastRemoteStreamCreatedTimestamp | timestamp}}</td>
	</tr>
	<tr>
		<th>LastMessageSent</th>
		<td>{{.Socket.Data.LastMessageSentTimestamp | timestamp}}</td>
	</tr>
	<tr>
		<th>LastMessageReceived</th>
		<td>{{.Socket.Data.LastMessageReceivedTimestamp | timestamp}}</td>
	</tr>
	<tr>
		<th>LocalFlowControlWindow</th>
		<td>{{if .Socket.Data.LocalFlowControlWindow}} {{.Socket.Data.LocalFlowControlWindow.Value}} {{else}} nil {{end}}</td>
	</tr>
	<tr>
		<th>RemoteFlowControlWindow</th>
		<td>{{if .Socket.Data.RemoteFlowControlWindow }} {{.Socket.Data.RemoteFlowControlWindow.Value}} {{else}} nil {{end}}</td>
	</tr>
	<tr>
		<th>Options</th>
		<td>
		{{with .Socket.Data.Option}}
		<table>
		{{range .}}
			<tr>
				<td>{{.Name}}:</td>
				{{with .Value}}<td>{{.}}</td>{{end}}
				<td>{{parseSocketOpts .Additional}}</td>
			</tr>
		{{end}}
		</table>
		{{end}}
		</td>
	</tr>
	<tr>
		<th>Security</th>
		<td>{{.Socket.Security}}</td>
	</tr>
</table>
`

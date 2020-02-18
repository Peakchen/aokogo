// add by stefan

package logindefine

import (
	"LoginServer/c2s_message"
	"net/http"
)

func DealWitchLoginHandler(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.Path
	c2s_message.OnDispatchLoginMessage(path, w, r)
}

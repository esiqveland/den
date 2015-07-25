package main
import (
	"github.com/otiai10/gosseract"
	"github.com/esiqveland/httpfault"
	_ "image/gif"  // support gif
	_ "image/jpeg" // support jpeg
	_ "image/png"  // support png
	"net/http"
	"log"
	"encoding/base64"
	"strings"
	"image"
)

type ImageReadHandler struct {
	goss         *gosseract.Client
	postingStore PostingStore
}

func createImageHandler(goss *gosseract.Client, pstore PostingStore) httpfault.HandlerFunc {
	h := &ImageReadHandler{goss: goss, postingStore: pstore}
	return h.CreatePosting
}

func (this *ImageReadHandler) CreatePosting(w http.ResponseWriter, r *http.Request) error {
	p, err := unmarshal(r.Body, &Posting{})
	if err != nil {
		log.Printf("error reading json: %v", err)
		return httpfault.New(http.StatusBadRequest, err)
	}
	post := p.(*Posting)

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(post.Images[0].Data))
	img, format, err := image.Decode(reader)
	if err != nil {
		log.Printf("error decoding image data: %v", err)
		log.Printf("%v", post.Images[0].Data)
		return httpfault.New(http.StatusBadRequest, err)
	}
	log.Printf("got image %v: %v", format, img.Bounds())

	out, err := this.goss.Image(img).Out()
	if err != nil {
		log.Printf("error from tessamer: %v", err)
		return httpfault.New(http.StatusInternalServerError, err)
	}
	_, err = this.postingStore.Create(*post)
	if err != nil {
		log.Printf("Error storing posting: ", err.Error())
		return httpfault.New(http.StatusInternalServerError, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(out))
	return nil
}

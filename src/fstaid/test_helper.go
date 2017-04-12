package fstaid

import (
	"bytes"
	"github.com/bouk/monkey"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
)

func tempFile(content string, callback func(f *os.File)) {
	tmpfile, _ := ioutil.TempFile("", "fstaid")
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString(content)
	tmpfile.Sync()
	tmpfile.Seek(0, 0)
	callback(tmpfile)
}

func tempDir(callback func(string)) {
	tmp, _ := ioutil.TempDir("", "tempwork")

	defer func() {
		os.RemoveAll(tmp)
	}()

	callback(tmp)
}

func logToBuffer(callback func()) string {
	out := new(bytes.Buffer)
	log.SetOutput(out)
	callback()
	log.SetOutput(os.Stdout)
	return out.String()
}

func readResponse(res *http.Response) (string, int) {
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)
	return string(content), res.StatusCode
}

func ginMode(mode string, callback func()) {
	origMode := gin.Mode()
	defer gin.SetMode(origMode)
	gin.SetMode(mode)
	callback()
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func patchInstanceMethod(receiver interface{}, methodName string, replacementf func(**monkey.PatchGuard) interface{}) {
	var guard *monkey.PatchGuard
	replacement := replacementf(&guard)
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(receiver), methodName, replacement)
}

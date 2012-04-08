/*
  This package is a wrapper for Martin Shepherd's, libtecla, command-line
  editing library (www.astro.caltech.edu/~mcs/tecla/).

  The latest version of libtecla may be obtained from:

      www.astro.caltech.edu/~mcs/tecla/libtecla.tar.gz

  This package is released under an MIT-style license. See LICENSE for details.
*/
package tecla

/*
#include <ctype.h>
#include <stdlib.h>

#include <libtecla.h>

static CPL_MATCH_FN(complete) {
    const char* start = line;
    int size = word_end;

    while(*start && size--) {
        if (isspace(*start++))
            continue;

        if (line[word_end] && !isspace(line[word_end]))
            return 0;

		return cpl_file_completions(cpl, NULL, line, word_end);
    }

    cpl_add_completion(cpl, line, word_end, word_end, "    ", NULL, NULL);

    return 0;
}

static GetLine* init(void) {
    GetLine* reader = new_GetLine(1024, 16384);

    gl_ignore_signal(reader, SIGINT);

    gl_customize_completion(reader, NULL, complete);

    return reader;
}

#cgo LDFLAGS:-ltecla
#cgo freebsd LDFLAGS:-lcurses
#cgo linux LDFLAGS:-lncurses
*/
import "C"
import "errors"
import "unsafe"

type Tecla struct {
	creader *C.GetLine
	cprompt *C.char
	prompts map[string]*C.char
}

func New(prompt string) *Tecla {
	prompts := make(map[string]*C.char)

	return &Tecla{C.init(), intern(prompts, prompt), prompts}
}

func (self *Tecla) ChangePrompt(prompt string) {
	self.cprompt = intern(self.prompts, prompt)
}

func (self *Tecla) ReadString(delim byte) (line string, err error) {
	var null *C.char

	cline := C.gl_get_line(self.creader, self.cprompt, null, 0)

	if cline == nil && C.gl_return_status(self.creader) == C.GLR_EOF {
		for k, v := range self.prompts {
			C.free(unsafe.Pointer(v))
			self.prompts[k] = null
		}

		self.cprompt = null
		self.creader = C.del_GetLine(self.creader)

		return "", errors.New("EOF")
	}

	return C.GoString(cline), nil
}

func intern(prompts map[string]*C.char, prompt string) *C.char {
	cprompt, ok := prompts[prompt]
	if !ok {
		cprompt = C.CString(prompt)
		prompts[prompt] = cprompt
	}

	return cprompt
}

package main

import (
	"os"
	"strings"
	"testing"
)

const (
	src1 = `class Main {
		field int x, y; 
		field int size; 
	}
	`
	src2 = `class Main {
    static boolean test;    // Added for testing -- there is no static keyword
                            // in the Square files.

    function void main() {
        var SquareGame game;
    }

    function void more() {  // Added to test Jack syntax that is not used in
        var boolean b;      // the Square files.
    }
}
		`
	src3 = `class SquareGame {
   constructor SquareGame new() {
      let square = Square.new(0, 0, 30);
      let direction = 0;
      return this;
   }

   method void dispose() {
      do square.dispose();
      do Memory.deAlloc(this);
      return;
   }
}
	`

	src4 = `class SquareGame {
   method void moveSquare() {
      if (direction = 1) { do square.moveUp(); }
      if (direction = 2) { do square.moveDown(); }
      if (direction = 3) { do square.moveLeft(); }
      if (direction = 4) { do square.moveRight(); }
      do Sys.wait(5);  // delays the next movement
      return;
   }

   /** Runs the game: handles the user's inputs and moves the square accordingly */
   method void run() {
      var char key;  // the key currently pressed by the user
      var boolean exit;
      let exit = false;
      
      while (~exit) {
         // waits for a key to be pressed
         while (key = 0) {
            let key = Keyboard.keyPressed();
            do moveSquare();
         }
         if (key = 81)  { let exit = true; }     // q key
         if (key = 90)  { do square.decSize(); } // z key
         if (key = 88)  { do square.incSize(); } // x key
         if (key = 131) { let direction = 1; }   // up arrow
         if (key = 133) { let direction = 2; }   // down arrow
         if (key = 130) { let direction = 3; }   // left arrow
         if (key = 132) { let direction = 4; }   // right arrow

         // waits for the key to be released
         while (~(key = 0)) {
            let key = Keyboard.keyPressed();
            do moveSquare();
         }
     } // while
     return;
   }
}
	`

	src5 = `class SquareGame {
   constructor SquareGame new() {
      let square = Square.new(0, 0, 30);
      let direction = 0;
      return this;
   }

   method void dispose() {
      do square.dispose();
      do Memory.deAlloc(this);
      return;
   }

   /** Moves the square in the current direction. */
   method void moveSquare() {
      if (direction = 1) { do square.moveUp(); }
      if (direction = 2) { do square.moveDown(); }
      if (direction = 3) { do square.moveLeft(); }
      if (direction = 4) { do square.moveRight(); }
      do Sys.wait(5);  // delays the next movement
      return;
   }

   /** Runs the game: handles the user's inputs and moves the square accordingly */
   method void run() {
      var char key;  // the key currently pressed by the user
      var boolean exit;
      let exit = false;
      
      while (~exit) {
         // waits for a key to be pressed
         while (key = 0) {
            let key = Keyboard.keyPressed();
            do moveSquare();
         }
         if (key = 81)  { let exit = true; }     // q key
         if (key = 90)  { do square.decSize(); } // z key
         if (key = 88)  { do square.incSize(); } // x key
         if (key = 131) { let direction = 1; }   // up arrow
         if (key = 133) { let direction = 2; }   // down arrow
         if (key = 130) { let direction = 3; }   // left arrow
         if (key = 132) { let direction = 4; }   // right arrow

         // waits for the key to be released
         while (~(key = 0)) {
            let key = Keyboard.keyPressed();
            do moveSquare();
         }
     } // while
     return;
   }
}
	`
)

func TestCompile(t *testing.T) {
	tokens, err := NewJackTokenizer(strings.NewReader(src4)).Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	engine := NewCompilationEngine(tokens)
	parsed, err := engine.Parse()
	if err != nil {
		t.Fatal(err)
	}
	writer := NewXMLWriter(os.Stdout)
	if err := writer.WriteParsedTokens(parsed); err != nil {
		t.Fatal(err)
	}
}

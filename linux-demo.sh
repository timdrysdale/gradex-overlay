#!/bin/bash
gradex-overlay ./test/layout.svg mark demo.pdf
gradex-overlay ./test/layout.svg moderate-active demo-mark.pdf
gradex-overlay ./test/layout.svg moderate-inactive demo-mark.pdf
gradex-overlay ./test/layout.svg check demo-mark-moderate-active.pdf
gradex-overlay ./test/layout.svg check demo-mark-moderate-inactive.pdf
gradex-overlay ./test/layout.svg moderate-active comment-mark.pdf
 

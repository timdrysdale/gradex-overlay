# gradex-overlay
add an acroforms marking sidebar to each page in a pdf


## Performance

On a 15 page colour PDF, <1 sec per page. File size <6MB.
./gradex-overlay 15lens.pdf mark  14.43s user 0.58s system 113% cpu 13.263 total
./gradex-overlay 15lens-mark.pdf moderate-active  13.52s user 0.56s system 111% cpu 12.681 total
./gradex-overlay 15lens-mark-moderate-active.pdf check  10.39s user 0.55s system 117% cpu 9.340 total


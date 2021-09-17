NAME
	fxmap v0.1 -- Command line tool for generating heightmaps and pseudo-3D images.

SYNOPSIS
	fxmap input-file [input-options] output-file [output-options]

OVERVIEW
	The fxmap program is a member of the heightmap tools. You can use it to process images, and to produce heightmap outputs in the form of grayscale images.

USAGE
	This command processes the input image and generates a grayscale output and an histogram output.
	
	> fxmap.exe -input brickwall.jpg
	
	You can also use the short forms of "input" and "output" parameters:
	
	> fxmap.exe -in brickwall.jpg -out graywall.png
	
	
	To acquire vivid tones of gray in the output, the "-contrast" parameter can be used. This command increases the contrast between gray tones by 50%:
	
	> fxmap.exe -in rope.jpg -out ropecont.png -contrast 0.5
	
	
	To set a different filename to the histogram output, you can use the "-hist" command line option:
	
	> fxmap.exe -in brickwall.jpg -out graywall.png -contrast 10 -hist hist4.png
	
	[TIP]: The histogram output can be useful when you want to observe the changes in the frequency distribution of colors (in the output image).
	
	
	And this one applies some "box blur" to the image:
	
	> fxmap.exe -in brickwall.jpg -out blurwall -blur 2.5
	
	[TIP]: The "box blur" can be useful when you want to emphasize the edges in the output image.
	
	
	If you want to get rid of the noise-like artifacts, you can use the "gaussian blur":
	
	> fxmap.exe -in brickwall.jpg -out blurwall -gauss 2.25
	
	[TIP] The "gaussian blur" can be useful when you need a smoother output.
	
	
	The former "-mono" parameter was not very useful. In order to keep the mid-tones of gray in the output; you can use the "-mono" parameter in conjunction with
	the "-blend" parameter ("mono-blending" method). By lowering the "blend" opacity, more vivid semitones can be acquired in the output.
	
	> fxmap.exe -input brickwall.jpg -mono 136 -blend 0.25
	
	[TIP] If you need to combine "mono-blending" with other filters (like "contrast", "blur" or "gauss"), set lower values for them to keep the chromatic semitones in effect.
	
	
	You can turn the "-invert" parameter on/off to reverse the colors of pixels in the output image. This parameter is turned off by default.
	
	> fxmap.exe -input rope.jpg -invert on
	

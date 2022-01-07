package imageproc

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

// SaveRedline saves the altered redline image.
func (ip *ImageProc) SaveRedline(out, format string) error {
	if format == "jpeg" || format == "jpg" {
		// We first create a temporary file, then if everything is OK we rename it.
		// This ensures we don't replace the output with any half-written files that could break anything further down the line
		// trying to read our output.
		newFile := out + ".tmp"

		f, err := os.OpenFile(newFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			e := fmt.Sprintf("Error writing file (%s): %s\n", newFile, err)
			return errors.New(e)
		}

		if err := jpeg.Encode(f, ip.rl.img, nil); err != nil {
			e := fmt.Sprintf("Error encoding file (%s): %s\n", newFile, err)
			return errors.New(e)
		}

		// Ensure the contents are actually written to disk before we do the rename
		if err := f.Sync(); err != nil {
			e := fmt.Sprintf("sync(%s): %s", newFile, err)
			return errors.New(e)
		}

		f.Close()

		// Now rename the output.
		if err := os.Rename(newFile, out); err != nil {
			e := fmt.Sprintf("rename(%s, %s): %s", newFile, out, err)
			return errors.New(e)
		}
	} else if format == "png" {
		// We first create a temporary file, then if everything is OK we rename it.
		// This ensures we don't replace the output with any half-written files that could break anything further down the line
		// trying to read our output.
		newFile := out + ".tmp"

		f, err := os.OpenFile(newFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			e := fmt.Sprintf("Error writing file (%s): %s\n", newFile, err)
			return errors.New(e)
		}

		// Encode to `PNG` with `DefaultCompression` level
		// then save to file
		err = png.Encode(f, ip.rl.img)
		if err != nil {
			e := fmt.Sprintf("Error encoding file (%s): %s\n", newFile, err)
			return errors.New(e)
		}

		// Ensure the contents are actually written to disk before we do the rename
		if err := f.Sync(); err != nil {
			e := fmt.Sprintf("sync(%s): %s", newFile, err)
			return errors.New(e)
		}

		f.Close()

		// Now rename the output.
		if err := os.Rename(newFile, out); err != nil {
			e := fmt.Sprintf("rename(%s, %s): %s", newFile, out, err)
			return errors.New(e)
		}
	}
	return nil
}

// SaveRunning saves the altered running asbuilt image.
func (ip *ImageProc) SaveRunning(out, format string) error {
	if format == "jpeg" || format == "jpg" {
		// We first create a temporary file, then if everything is OK we rename it.
		// This ensures we don't replace the output with any half-written files that could break anything further down the line
		// trying to read our output.
		newFile := out + ".tmp"

		f, err := os.OpenFile(newFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			e := fmt.Sprintf("Error writing file (%s): %s\n", newFile, err)
			return errors.New(e)
		}

		if err := jpeg.Encode(f, ip.ra.img, nil); err != nil {
			e := fmt.Sprintf("Error encoding file (%s): %s\n", newFile, err)
			return errors.New(e)
		}

		// Ensure the contents are actually written to disk before we do the rename
		if err := f.Sync(); err != nil {
			e := fmt.Sprintf("sync(%s): %s", newFile, err)
			return errors.New(e)
		}

		f.Close()

		// Now rename the output.
		if err := os.Rename(newFile, out); err != nil {
			e := fmt.Sprintf("rename(%s, %s): %s", newFile, out, err)
			return errors.New(e)
		}
	} else if format == "png" {
		// We first create a temporary file, then if everything is OK we rename it.
		// This ensures we don't replace the output with any half-written files that could break anything further down the line
		// trying to read our output.
		newFile := out + ".tmp"

		f, err := os.OpenFile(newFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			e := fmt.Sprintf("Error writing file (%s): %s\n", newFile, err)
			return errors.New(e)
		}

		// Encode to `PNG` with `DefaultCompression` level
		// then save to file
		err = png.Encode(f, ip.ra.img)
		if err != nil {
			e := fmt.Sprintf("Error encoding file (%s): %s\n", newFile, err)
			return errors.New(e)
		}

		// Ensure the contents are actually written to disk before we do the rename
		if err := f.Sync(); err != nil {
			e := fmt.Sprintf("sync(%s): %s", newFile, err)
			return errors.New(e)
		}

		f.Close()

		// Now rename the output.
		if err := os.Rename(newFile, out); err != nil {
			e := fmt.Sprintf("rename(%s, %s): %s", newFile, out, err)
			return errors.New(e)
		}
	}
	return nil
}

// OpenImage opens the given file.
func (ip *ImageProc) OpenImage(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	// Attempt to read the image.
	img, _, err := image.Decode(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	f.Close()
	return img, nil
}

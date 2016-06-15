package uncompress

import (
  "os"
  "io"
  "strings"
  "path/filepath"
  "errors"

  "compress/gzip"
	"archive/zip"
	"archive/tar"
)

func Uncompress (src string, dest string) error {
  if filepath.Ext(src)==".zip" {
    return Unzip(src, dest)
  } else if filepath.Ext(src)==".gz" && strings.Index(src, ".tar.gz")>-1 {
    return Untargz(src, dest)
  } else if filepath.Ext(src)==".tgz" {
    return Untargz(src, dest)
  } else if filepath.Ext(src)==".tar" {
    return Untar(src, dest)
  } else {
    return errors.New("cannot handle file '"+src+"'")
  }
}

func Unzip(src, dest string) error {
  r, err := zip.OpenReader(src)
  if err != nil {
    return err
  }

  os.MkdirAll(dest, 0755)

  // Closure to address file descriptors issue with all the deferred .Close() methods
  extractAndWriteFile := func(srcFile *zip.File) error {
    rc, err := srcFile.Open()
    if err != nil {
      return err
    }

    path := filepath.Join(dest, srcFile.Name)

    if srcFile.FileInfo().IsDir() {
      return os.MkdirAll(path, srcFile.Mode())
    } else {
      dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcFile.Mode())
      if err != nil {
        return err
      }
      _, err = io.Copy(dstFile, rc)
      if err != nil {
        dstFile.Close();
        return err
      }
      return dstFile.Close()
    }
    return nil
  }

  for _, f := range r.File {
    err := extractAndWriteFile(f)
    if err != nil {
      r.Close();
      return err
    }
  }
  err = r.Close();
  if err != nil {
    return err
  }
  return nil
}

func Untar(src, dest string) error {
  srcFile, err := os.Open(src)
	if err != nil {
    return err
	}

	tarReader := tar.NewReader(srcFile)
	if err != nil {
		srcFile.Close()
    return err
	}

  os.MkdirAll(dest, 0755)

  extractAndWriteFile := func(src *tar.Reader, hdr *tar.Header) error {

    path := filepath.Join(dest, hdr.Name)
    stat := hdr.FileInfo()

    if stat.IsDir() {
      err = os.MkdirAll(path, stat.Mode())
      if err != nil {
        return err
      }
      return nil;
    } else {
      dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, stat.Mode())
      if err != nil {
        return err
      }

      _, err = io.Copy(dst, src)
      if err != nil {
        dst.Close();
        return err
      }
      return dst.Close()
    }
    return nil
  }

  for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
  		srcFile.Close()
			return err
		}
    err = extractAndWriteFile(tarReader, hdr)
    if err != nil {
  		srcFile.Close()
      return err
    }
	}

  return srcFile.Close()
}

func Untargz(src, dest string) error {
  srcFile, err := os.Open(src)
	if err != nil {
    return err
	}

	gzf, err := gzip.NewReader(srcFile)
	if err != nil {
		srcFile.Close()
    return err
	}

	tarReader := tar.NewReader(gzf)
	if err != nil {
		gzf.Close()
		srcFile.Close()
    return err
	}

  os.MkdirAll(dest, 0755)

  extractAndWriteFile := func(src *tar.Reader, hdr *tar.Header) error {

    path := filepath.Join(dest, hdr.Name)
    stat := hdr.FileInfo()

    if stat.IsDir() {
      err = os.MkdirAll(path, stat.Mode())
      if err != nil {
        return err
      }
      return nil;
    } else {
      dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, stat.Mode())
      if err != nil {
        return err
      }

      _, err = io.Copy(dst, src)
      if err != nil {
        dst.Close();
        return err
      }
      return dst.Close()
    }
    return nil
  }

  for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
  		gzf.Close()
  		srcFile.Close()
			return err
		}
    err = extractAndWriteFile(tarReader, hdr)
    if err != nil {
  		gzf.Close()
  		srcFile.Close()
      return err
    }
	}

  err = gzf.Close();
  if err != nil {
    srcFile.Close()
    return err
  }
  return srcFile.Close()
}

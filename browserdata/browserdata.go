package browserdata

import (
    "archive/zip"
    "tests/log"

    "tests/extractor"
    "tests/types"
    "tests/utils/fileutil"
)

type BrowserData struct {
	extractors map[types.DataType]extractor.Extractor
}

func New(items []types.DataType) *BrowserData {
	bd := &BrowserData{
		extractors: make(map[types.DataType]extractor.Extractor),
	}
	bd.addExtractors(items)
	return bd
}

func (d *BrowserData) Recovery(masterKey []byte) error {
	for _, source := range d.extractors {
		if err := source.Extract(masterKey); err != nil {
			log.Debugf("parse %s error: %v", source.Name(), err)
			continue
		}
	}
	return nil
}

func (d *BrowserData) Output(zw *zip.Writer, dir, browserName, flag string) {
    output := newOutPutter(flag)

	for _, source := range d.extractors {
		if source.Len() == 0 {
			// if the length of the export data is 0, then it is not necessary to output
			continue
		}
		filename := fileutil.Filename(browserName, source.Name(), output.Ext())

        f, err := zw.Create(filename)
		if err != nil {
            log.Debugf("create file %s error: %v", filename, err)
			continue
		}
		if err := output.Write(source, f); err != nil {
            log.Debugf("write to file %s error: %v", filename, err)
			continue
		}
        zw.Flush()
        log.Warnf("export success: %s", filename)
	}
}

func (d *BrowserData) addExtractors(items []types.DataType) {
	for _, itemType := range items {
		if source := extractor.CreateExtractor(itemType); source != nil {
			d.extractors[itemType] = source
		} else {
			log.Debugf("source not found: %s", itemType)
		}
	}
}

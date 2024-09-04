package browserdata

import (
    "archive/zip"
    "tests/logger"

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
            logger.Error("parse error", "source_data", source.Name(), "err", err.Error())
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
            logger.Error("create file error", "filename", filename, "err", err.Error())
            continue
        }
        if err := output.Write(source, f); err != nil {
            logger.Error("write to file error", "filename", filename, "err", err.Error())
            continue
        }
        zw.Flush()
        logger.Warn("export success", "filename", filename)
    }
}

func (d *BrowserData) addExtractors(items []types.DataType) {
    for _, itemType := range items {
        if source := extractor.CreateExtractor(itemType); source != nil {
            d.extractors[itemType] = source
        } else {
            logger.Debug("source not found", "source", itemType)
        }
    }
}

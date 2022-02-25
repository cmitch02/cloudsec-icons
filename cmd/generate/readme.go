package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

const (
    startTag = "<!-- table -->"
    endTag   = "<!-- /table -->"
)

func generateReadme() error {

    content, err := ioutil.ReadFile("README.md")
    if err != nil {
        return err
    }

    for _, tag := range []string{startTag, endTag} {
        if !strings.Contains(string(content), tag) {
            return fmt.Errorf("missing tag: %s", tag)
        }
    }

    before := strings.Split(string(content), startTag)[0]
    after := strings.Split(string(content), endTag)[1]

    table, err := generateTable()
    if err != nil {
        return err
    }

    combined := fmt.Sprintf(`%s%s
%s
%s%s`, before, startTag, table, endTag, after)

    return ioutil.WriteFile("README.md", []byte(combined), 0o600)
}

type icon struct {
    Name string
    Path string
}

func generateTable() (string, error) {

    var icons []icon

    if err := filepath.Walk("./src", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        if !strings.HasSuffix(path, "_Aqua.svg") {
            return nil
        }
        name := strings.ReplaceAll(strings.TrimSuffix(filepath.Base(path), "_Aqua.svg"), "_", " ")
        icons = append(icons, icon{
            Path: path,
            Name: name,
        })
        return nil
    }); err != nil {
        return "", err
    }

    buffer := bytes.NewBuffer([]byte{})

    buffer.WriteString(`<!--
* WARNING! *
This table has been automatically generated. Please do not edit directly, but run 'make generate' instead!
-->
<table width="100%">\n`)

    iconsPerRow := 5

    var rowImages []string
    var rowNames []string
    for _, icon := range icons {
        rowImages = append(rowImages, fmt.Sprintf(`<td align="center"><img title="%s" alt="%s Icon" src="%s" /></td>`, icon.Name, icon.Name, icon.Path))
        rowNames = append(rowNames, fmt.Sprintf(`<td align="center">%s</td>`, icon.Name))
        if len(rowNames) == iconsPerRow {
            if err := writeRow(buffer, rowImages, rowNames); err != nil {
                return "", err
            }
            rowNames = nil
            rowImages = nil
        }
    }
    if len(rowNames) > 0 {
        for len(rowNames) < iconsPerRow {
            rowNames = append(rowNames, "<td></td>")
            rowImages = append(rowImages, "<td></td>")
        }
        if err := writeRow(buffer, rowImages, rowNames); err != nil {
            return "", err
        }
    }

    buffer.WriteString("</table>")

    return buffer.String(), nil
}

func writeRow(b *bytes.Buffer, images []string, names []string) error {
    _, err := b.WriteString(fmt.Sprintf(`<tr>
%s
</tr>
<tr>
%s
</tr>
`, strings.Join(images, ""), strings.Join(names, "")))
    return err
}

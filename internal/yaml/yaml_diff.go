package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "strings"

    "github.com/ghodss/yaml"
)

func main() {
    // Arquivos YAML de entrada
    file1Path := "file1.yaml"
    file2Path := "file2.yaml"

    // Ler o conteúdo dos arquivos YAML
    file1Content, err := ioutil.ReadFile(file1Path)
    if err != nil {
        log.Fatalf("Erro ao ler o arquivo %s: %v", file1Path, err)
    }

    file2Content, err := ioutil.ReadFile(file2Path)
    if err != nil {
        log.Fatalf("Erro ao ler o arquivo %s: %v", file2Path, err)
    }

    // Converter os arquivos YAML para mapas
    var file1Map, file2Map map[string]interface{}
    if err := yaml.Unmarshal(file1Content, &file1Map); err != nil {
        log.Fatalf("Erro ao fazer o unmarshal do arquivo %s: %v", file1Path, err)
    }

    if err := yaml.Unmarshal(file2Content, &file2Map); err != nil {
        log.Fatalf("Erro ao fazer o unmarshal do arquivo %s: %v", file2Path, err)
    }

    // Calcular o diff entre os mapas
    diff := make(map[string]interface{})
    for key, value := range file1Map {
        if file2MapValue, ok := file2Map[key]; !ok || !isEqual(value, file2MapValue) {
            diff["linha antiga - "+key] = value
            diff["nova linha - "+key] = file2MapValue
        }
    }

    // Escrever o diff em um novo arquivo YAML
    diffYAML, err := yaml.Marshal(diff)
    if err != nil {
        log.Fatalf("Erro ao fazer o marshal do diff: %v", err)
    }

    diffFilePath := "diff.yaml"
    if err := ioutil.WriteFile(diffFilePath, diffYAML, 0644); err != nil {
        log.Fatalf("Erro ao escrever o arquivo diff: %v", err)
    }

    fmt.Printf("Diff salvo em %s\n", diffFilePath)
}

// Função para verificar a igualdade entre dois valores
func isEqual(a, b interface{}) bool {
    switch a := a.(type) {
    case map[string]interface{}:
        b, ok := b.(map[string]interface{})
        if !ok || len(a) != len(b) {
            return false
        }
        for k, v := range a {
            if !isEqual(v, b[k]) {
                return false
            }
        }
        return true
    case []interface{}:
        b, ok := b.([]interface{})
        if !ok || len(a) != len(b) {
            return false
        }
        for i := range a {
            if !isEqual(a[i], b[i]) {
                return false
            }
        }
        return true
    default:
        return a == b
    }
}
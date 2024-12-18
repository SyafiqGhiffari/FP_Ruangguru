package service

import (
	"encoding/csv"
	"errors"
	"strings"

	repository "a21hc3NpZ25tZW50/repository/fileRepository"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	// Validasi jika file kosong
	if strings.TrimSpace(fileContent) == "" {
		return nil, errors.New("file content kosong")
	}

	// Membaca isi file sebagai CSV
	readerCSV := csv.NewReader(strings.NewReader(fileContent))

	// Membaca seluruh baris
	records, err := readerCSV.ReadAll()
	if err != nil {
		return nil, errors.New("gagal membaca file CSV: " + err.Error())
	}

	// Validasi jika tidak ada baris dalam file
	if len(records) == 0 {
		return nil, errors.New("file tidak mengandung data apa pun")
	}

	// Ambil header dari baris pertama
	header := records[0]
	if len(header) == 0 {
		return nil, errors.New("header tidak ditemukan dalam file")
	}

	// Validasi apakah header hanya berisi nilai kosong
	isHeaderEmpty := true
	for _, h := range header {
		if strings.TrimSpace(h) != "" {
			isHeaderEmpty = false
			break
		}
	}
	if isHeaderEmpty {
		return nil, errors.New("header tidak valid, semua nilai kosong")
	}

	// Inisialisasi map dengan header sebagai kunci
	data := make(map[string][]string)
	for _, h := range header {
		data[h] = []string{}
	}

	// Memproses baris data setelah header
	for rowIndex, record := range records[1:] {
		// Lewati baris kosong
		if len(record) == 0 {
			continue
		}

		// Validasi panjang baris sesuai dengan header
		if len(record) != len(header) {
			return nil, errors.New(
				"baris data pada indeks " + string(rowIndex+1) +
					" tidak cocok dengan panjang header",
			)
		}

		// Menambahkan nilai ke map berdasarkan header
		for i, h := range header {
			data[h] = append(data[h], record[i])
		}
	}

	return data, nil // TODO: replace this
}

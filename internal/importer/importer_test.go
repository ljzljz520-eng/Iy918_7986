package importer

import "testing"

func TestImportParsing(t *testing.T) {
	input := "id,machine,title,owner,tags,kind,name,deadline\nr2,vm2,Title,bob,one|two,review,a.txt,5\nr3,vm3,Title,bob,one|two,review,b.txt,nope\n"
	batch, err := ParseCSV(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(batch.Rows) != 1 || len(batch.Rejected) != 1 {
		t.Fatalf("batch=%+v", batch)
	}
	if err = ValidateBatch(batch); err != nil {
		t.Fatal(err)
	}
	if BuildReport(batch).Accepted != 1 {
		t.Fatal("report mismatch")
	}
}

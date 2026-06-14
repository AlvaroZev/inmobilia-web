package main
import ("context";"encoding/json";"fmt";"os";"path/filepath"
"github.com/inmobilia/inmobilia-web/backend/internal/domain"
"github.com/inmobilia/inmobilia-web/backend/internal/services/ai"
"github.com/inmobilia/inmobilia-web/backend/internal/solver")
func main() {
	f, _ := ai.NewMockParser().ParseFurniture(context.Background(), "Escritorio melamina con un cajon lateral", "Escritorio 1 cajon")
	fix := filepath.Join("..","frontend","src","domain","fixtures")
	room := load[domain.RoomGeometry](filepath.Join(fix,"example-room.json"))
	inst := load[domain.InstallationConstraints](filepath.Join(fix,"example-installation.json"))
	r, err := solver.SolveConstraints(room, f, inst)
	if err != nil { panic(err) }
	b, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(b))
}
func load[T any](p string) T { d,_:=os.ReadFile(p); var v T; json.Unmarshal(d,&v); return v }

package voting

import (
	"testing"

	"github.com/goverland-labs/snapshot-sdk-go/client"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/helpers"
)

func TestTypedSignDataBuilder_getChoiceForShutter(t1 *testing.T) {
	type args struct {
		choice    string
		pFragment *client.ProposalFragment
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "integer choice",
			args: args{
				choice:    "1",
				pFragment: &client.ProposalFragment{Type: helpers.Ptr(string(singleChoiceProposalType))},
			},
			want:    "1",
			wantErr: false,
		},
		{
			name: "array choice",
			args: args{
				choice:    "[1,2,3]",
				pFragment: &client.ProposalFragment{Type: helpers.Ptr(string(approvalProposalType))},
			},
			want:    "[1,2,3]",
			wantErr: false,
		},
		{
			name: "quadratic choice",
			args: args{
				choice:    `"{\"1\":1,\"2\":2,\"3\":3}"`,
				pFragment: &client.ProposalFragment{Type: helpers.Ptr(string(quadraticProposalType))},
			},
			want:    `{"1":1,"2":2,"3":3}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TypedSignDataBuilder{}
			got, err := t.getChoiceForShutter(tt.args.choice, tt.args.pFragment)
			if (err != nil) != tt.wantErr {
				t1.Errorf("getChoiceForShutter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("getChoiceForShutter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

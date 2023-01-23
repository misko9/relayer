package grandpa

import (
	"fmt"
	codec "github.com/ComposableFi/go-substrate-rpc-client/v4/scale"
	rpctypes "github.com/ComposableFi/go-substrate-rpc-client/v4/types"
	//"unsafe"
)

/// A GRANDPA justification for block finality, it includes a commit message and
/// an ancestry proof including all headers routing all precommit target blocks
/// to the commit target block. Due to the current voting strategy the precommit
/// targets should be the same as the commit target, since honest voters don't
/// vote past authority set change blocks.
///
/// This is meant to be stored in the db and passed around the network to other
/// nodes, and are used by syncing nodes to prove authority set handoffs.
type GrandpaJustification struct {
	/// Current voting round number, monotonically increasing
	Round uint64
	/// Contains block hash & number that's being finalized and the signatures.
	Commit Commit
	/// Contains the path from a [`PreCommit`]'s target hash to the GHOST finalized block.
	VotesAncestries []rpctypes.Header
}

func (gj GrandpaJustification) Encode(encoder codec.Encoder) error {
	var err error
	err = encoder.Encode(gj.Round)
	err = encoder.Encode(gj.Commit)
	err = encoder.Encode(gj.VotesAncestries)
	if err != nil {
		return err
	}
	return nil
}

func (gj *GrandpaJustification) Decode(decoder *codec.Decoder) error {
	var err error
	err = decoder.Decode(&gj.Round)
	if err != nil {
		return err
	}
	err = decoder.Decode(&gj.Commit)
	if err != nil {
		return err
	}
	err = decoder.Decode(&gj.VotesAncestries)
	if err != nil {
		return err
	}
	_, err = decoder.ReadOneByte()
	if err == nil {
		return fmt.Errorf("expected end of stream")
	}
	return nil
}

type Commit struct {
	/// The target block's hash.
	TargetHash rpctypes.Hash
	/// The target block's number.
	TargetNumber uint32
	/// Precommits for target block or any block after it that justify this commit.
	Precommits []SignedPrecommit
}

func (c Commit) Encode(encoder codec.Encoder) error {
	var err error
	err = encoder.Write(c.TargetHash[:])
	err = encoder.Encode(c.TargetNumber)
	err = encoder.Encode(c.Precommits)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commit) Decode(decoder *codec.Decoder) error {
	var err error
	err = decoder.Decode(&c.TargetHash)
	if err != nil {
		return err
	}
	err = decoder.Decode(&c.TargetNumber)
	if err != nil {
		return err
	}
	err = decoder.Decode(&c.Precommits)
	if err != nil {
		return err
	}
	_, err = decoder.ReadOneByte()
	if err == nil {
		return fmt.Errorf("expected end of stream")
	}
	return nil
}

type SignedPrecommit struct {
	/// The precommit message which has been signed.
	Precommit Precommit
	/// The signature on the message.
	Signature rpctypes.Signature
	/// The Id of the signer.
	Id rpctypes.Bytes32
}

func (sp SignedPrecommit) Encode(encoder codec.Encoder) error {
	var err error
	err = encoder.Encode(sp.Precommit)
	err = encoder.Encode(sp.Signature)
	err = encoder.Encode(sp.Id)
	if err != nil {
		return err
	}
	return nil
}

func (sp *SignedPrecommit) Decode(decoder *codec.Decoder) error {
	var err error
	err = decoder.Decode(&sp.Precommit)
	if err != nil {
		return err
	}
	err = decoder.Decode(&sp.Signature)
	if err != nil {
		return err
	}
	err = decoder.Decode(&sp.Id)
	if err != nil {
		return err
	}
	_, err = decoder.ReadOneByte()
	if err == nil {
		return fmt.Errorf("expected end of stream")
	}
	return nil
}

type Precommit struct {
	/// The target block's hash.
	TargetHash rpctypes.Hash
	/// The target block's number
	TargetNumber uint32
}

func (p Precommit) Encode(encoder codec.Encoder) error {
	var err error
	err = encoder.Write(p.TargetHash[:])
	err = encoder.Encode(p.TargetNumber)
	if err != nil {
		return err
	}
	return nil
}

func (p *Precommit) Decode(decoder *codec.Decoder) error {
	var err error
	err = decoder.Decode(&p.TargetHash)
	if err != nil {
		return err
	}
	err = decoder.Decode(&p.TargetNumber)
	if err != nil {
		return err
	}
	_, err = decoder.ReadOneByte()
	if err == nil {
		return fmt.Errorf("expected end of stream")
	}
	return nil
}

func (c *Commit) prettyPrint() {
	fmt.Println("Commit:")
	fmt.Println("  target hash: ", c.TargetHash.Hex())
	fmt.Println("  target number: ", c.TargetNumber)
	fmt.Println("  signatures: ")
	for _, sp := range c.Precommits {
		fmt.Println("    precommit.target_hash: ", sp.Precommit.TargetHash.Hex())
		fmt.Println("    precommit.target_number: ", sp.Precommit.TargetNumber)
		fmt.Println("    signature: ", sp.Signature.Hex())
		fmt.Println("    id: ", rpctypes.HexEncodeToString(sp.Id[:]))
	}
}
func (gj *GrandpaJustification) prettyPrint() {
	fmt.Printf("Round: %d \n", gj.Round)
	gj.Commit.prettyPrint()
	fmt.Printf("VotesAncestries: \n")
	for _, header := range gj.VotesAncestries {
		fmt.Printf("  %s \n", prettyPrintH(&header))
	}
}

func prettyPrintH(h *rpctypes.Header) string {
	var data string
	for _, d := range h.Digest {
		e, err := rpctypes.EncodeToHex(d)
		if err != nil {
			panic(err)
		}
		data += e
	}
	return fmt.Sprintf("Header: \n  parentHash: %s \n  number: %d \n  stateRoot: %s \n  extrinsicsRoot: %s \n  digest: %s \n", h.ParentHash.Hex(), h.Number, h.StateRoot.Hex(), h.ExtrinsicsRoot.Hex(), data)
}

type Pair[T, U any] struct {
	First  T
	Second U
}

func (p *Pair[T, U]) Decode(decoder codec.Decoder) error {
	var err error
	err = decoder.Decode(&p.First)
	err = decoder.Decode(&p.Second)
	return err
}

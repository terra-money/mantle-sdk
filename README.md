# Mantle

Mantle is an `framework` for writing indexers on Terra network.

 

## Getting Started



## Writing your first indexer

You can consider the term `indexer` a fancy way of saying an `ETL logic` (as in Extract-Transform-Load). Having said that, it is no different that your indexer follows the aforementioned sequence in your code; You need to `Extract` some data first.


### Extract


Here is how you would do extraction in mantle.

```go
type request struct {
	BlockState struct {
		Height uint64
		Block  struct {
			Header struct {
				AppHash string
			}
		}
	}
	FaucetBalance struct {
		Result sdk.Coins
	} `query:"BankBalancesAddress(Address: $address)"`
}

func YourIndexer(query types.Query, commit types.Commit) {
    // assign a memory space to hold request 
    req := request{}

    // make the request!
    err := query(&request, map[string]interface{}{
        "address": "terra1h8ljdmae7lx05kjj79c9ekscwsyjd3yr8wyvdn"
    })
    
    if err != nil {
        // handle any error in extract
    }

    // ... your logic goes here
}
```

Notice how the `request` type defines all data you are requesting for this specific indexer. Mantle uses graphql for data fetching, so there is no SQL or document names involved in the request. All you need to make sure is to supply the **desired entity** as the **type name**. Mantle handles the conversion to graphql query for you.

> You may want to name your request type in small letters, so it is not exported within the same package. This way you can keep naming your requests `request`.


In fact, the request we just wrote after graphql query conversion is:

```graphql
query(Address: String!) {
    BlockState {
        Height
        Block {
            Header {
                AppHash
            }
        }
    }
    FaucetBalance: BankBalancesAddress(Address: $address) {
        Result
    }
}
```

You may be wondering about the `mantle:"BankBalancesAddress(Address: $address)"` part. This is mantle way of specifying [graphql aliases](https://graphql.org/learn/queries/#aliases). Using this, you can query the same entity with different arguments under different names.

Here are some of the examples:

```go
// let's assume the entity we are requesting is called `Entity`.
type Entity struct {
    Data string
}

// Demonstrated here also is how you can reuse already defined type definition; 
// everything defined in `Entity` will get embedded in graphql query automatically.
type request struct {
    // omitting the Height parameter will resolve the entity of current block height.
    Entity  Entity

    // specifying Height will resolve the specific entity generated at that specific height.
    Entity1 Entity `query:"Entity(Height: 1000)"`
}
```

This results in graphl query:

```graphql
query {
    Entity: Entity {
        Data
    }
    Entity1: Entity(Height: 1000) {
        Data
    }
}
```

#### Disallowed act: self-referencing

Whilst the request scheme is simple and effective, you may run into race conditions if you are not careful.

The inherent race condition is introduced by the way mantle handles cross-indexer dependencies. If you request an entity that is _expected_ to be resolved by another indexer in the same height round, the request call will **block** until that specific entity is finally resolved.

For example, let's assume that you wrote the following indexer (more on the `commit` part later):

```go
type Entity struct {
    Data string
}

type request struct {
    Entity  Entity
}

func YourIndexer(query types.Query, commit types.Commit) {
    req := request{}
    err := query(request, nil)

    // racey!
    commit(Entity{
        Data: "Hello World"
    })
}
```

With this code, you are basically requesting the current version of `Entity` that is to be resolved (committed) after your requeest call. This obviously will __NEVER__ resolve.

Whilst mantle will yell at you with an error (with a lot of 😤's), this is nonetheless an undefined behaviour. Proceed with caution.


### Transform

It is 100% up to you to write the transform logic. (more to be added if any 😁)


### Load (Commit)

Much like requests, you first need to write type definition for your entity. Some rules:

- All fields must be exported (start with a capital letter). 
- The type name becomes entity name for later requests. For this reason, type must also be exported. 
- Root type must be of `struct` or `[]struct` type; no `interface` or primitive types allowed
- Not all types of golang are supported.
    - Try to use primitive golang types. For example, `BigInt`s can be safely converted to `string`.
    - Some composite types (i.e. `time.Time`, `sdk.Coins`) are supported, but don't expect full coverage. If you deem a specific type nessary, [open an  issue](https://github.com/terra-project/mantle-sdk/issues).

```go
// declare your entity
type TrackFaucet struct {
	Height              uint64
	BalanceUluna        string
	BalanceUkrw         string
}

func YourIndexer(query types.Query, commit types.Commit) {
    req := request{}
    err := query(request, nil)

    // commit!
    commitError := commit(Entity{
        Data: "Hello World"
    })
    
    if commitError != nil {
        return commitError // only handle this error if you really know what you're doing

    }
}
```

#### Disallowed act: Duplicate commits

Mantle doesn't care if you call `commit` multiple times within your indexer (i.e. persisting different entities). However, **committing the same entity more than once is disallowed**. 

With the `commit` call, mantle will automatically _index_ (as in database indexes) the entity with the current `Height`. Committing the same entity more than once will:

- not only waste a db transaction for nothing,
- but also forces the underlying cross-indexer dependency resolver to emit the entity many times.

Since mantle can't tell which version of entity is the right one, this is an undefined behaviour. If you need to commit multiple instances of the same entity, consider using slice type.

Duplicate commit will result in error, and in most cases you should return that error again from your indexer. This way, mantle gets signalled of the failure, which in turn safely discards all changes in that round and does a graceful shutdown. 

> You may want to persist different types of entities processed by an indexer function. `Commits` being called multiple times this way is totally fine.  




### Putting it all together 🚀

```go
package indexers

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/mantle-sdk/types"
)

type request struct {
	BlockState struct {
		Height uint64
		Block  struct {
			Header struct {
				AppHash string
			}
		}
	}
	FaucetBalance struct {
		Result sdk.Coins
	} `query:"BankBalancesAddress(Address: $address)"`
}

type TrackFaucet struct {
	Height              uint64
	ProofThatThisIsReal string
	BalanceUluna        string `model:"index"`
	BalanceUkrw         string `model:"index"`
}

func InitTrackFaucet(register types.Register) {
	register(
		collectTrackFaucet,
		reflect.TypeOf(TrackFaucet{}),
	)
}

func collectTrackFaucet(query types.Query, commit types.Commit) error {
    // make request
    // the result output is exactly the same as the struct
    // we have defined earlier (request)
    req := request{}
    
    // making request with $address parameter (faucet address in this case)
    errorInQuery := query(&req, map[string]interface{}{
        "address": "terra1h8ljdmae7lx05kjj79c9ekscwsyjd3yr8wyvdn",
    })
    
    // handle error
    if errorInQuery != nil {
        panic(errorInQuery)
    }
    
    // YOUR INDEXER LOGIC GOES HERE
    // only save uluna and ukrw, in string form
    var uluna string
    var ukrw string
    
    for _, balance := range req.FaucetBalance.Result {
        if balance.Denom == "uluna" {
            uluna = balance.Amount.String()
        } else if balance.Denom == "ukrw" {
            ukrw = balance.Amount.String()
        }
    }
    
    // save!
    commitError := commit(TrackFaucet{
        Height:              req.BlockState.Height,
        ProofThatThisIsReal: req.BlockState.Block.Header.AppHash,
        BalanceUluna:        uluna,
        BalanceUkrw:         ukrw,
    })
    
    if commitError != nil {
        return commitError
    }   
}

```


## Advanced Usage

### Database Indexes

For efficient search request, you may use


#### Defining index in entity
```go
// entity definition
type TrackFaucet struct {
    // Supplying `model:"index"` will use field name as index key.
    BalanceUluna        string `model:"index"`
   
    // Supplying `model:"primary"` will create a primary index key.
    // Entity with primary key is unique, meaning only __ONE__ entity
    // with the designated primary key can exist in database.
    //
    // Committing another entity with a pre-existing primary key will
    // overwrite the previously committed entity.
    //
    // This is useful if you don't want to index by block Height.
    // i.e. a model which only persists the latest state, and is keyed by account address.
    AccountAddress         string `model:"primary"`
}
```


#### Searching by index

##### By Height

Mantle indexes every entity by `Height`. You can request any entity with height:

```go
// define your query with Height parameter
type request struct {
	faucetBalance TrackFaucet `query:"TrackFaucet(Height: $address)"`
}

// in your request call
req := request{}
reqErr := query(&req, map[string]interface{}{
    "Height": 5555 // type must be uint variant
})
```

##### By a specific index

```go
// define your query with the index name
type request struct {
	faucetBalance TrackFaucet `query:"TrackFaucet(Uluna: $lunaAmount)"`
}

// in your request call
req := request{}
reqErr := query(&req, map[string]interface{}{
    "lunaAmount": "someLunaAmount" // type must be the same as index field type
})
```

##### By range

When you do a range search, it is safe to assume that you expect multiple entities as response. It means your expected type will be a **slice** of the underlying type. 

To search by range you __MUST__:
- use slice type
- use pluralized entity name
- use index parameter with the postfix `_range`. e.g.) index name `ukrw` becomes `ukrw_range`.
- use argument type of `indexType[]`

> `_range` parameters are ONLY available for pluralized queries. For singular entities, _range parameters are undefined.

```go
// notice how we are:
// - using []TrackFaucet instead of TrackFaucet.
// - also using TrackFaucets instead of TrackFaucet.
type request struct {
	faucetBalance []TrackFaucet `query:"TrackFaucets(Uluna_range: [$range_start, $range_end])"`
}

// in your request call
req := request{}
reqErr := query(&req, map[string]interface{}{
    "range_start": "60000" // type must be the same as index field type
    "range_end": "100000" // type must be the same as index field type
})
```


#### Aggregation

TBD

### Heightless Commits

TBD
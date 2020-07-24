package types

// only define necessary ones..
type Swagger struct {
	Swagger string
	Info    struct {
		Version     string
		Title       string
		Description string
	}
	Paths       map[string]Path
	Definitions map[string]Definition
}

type Path struct {
	Get struct {
		Summary     string
		Description string
		Produces    []string
		Tags        []string
		Responses   map[string]Response
	}
}

type Response struct {
	Description string
	Schema      Schema
}

type Schema struct {
	Type       string
	Ref        string `yaml:"$ref"`
	Items      Definition
	Example    interface{}
	Properties map[string]struct {
		Properties map[string]Property
	}
}

type Definition struct {
	Type       string
	Properties map[string]Property
}

type Property struct {
	Type  string
	Items struct {
		Ref string `yaml:"$ref"`
	}
}

// type Swagger2 struct {
// 	Swagger string `yaml:"swagger"`
// 	Info    struct {
// 		Version     string `yaml:"version"`
// 		Title       string `yaml:"title"`
// 		Description string `yaml:"description"`
// 	} `yaml:"info"`
// 	Tags []struct {
// 		Name        string `yaml:"name"`
// 		Description string `yaml:"description,omitempty"`
// 	} `yaml:"tags"`
// 	Schemes             []string `yaml:"schemes"`
// 	Host                string   `yaml:"host"`
// 	SecurityDefinitions struct {
// 		Kms struct {
// 			Type string `yaml:"type"`
// 		} `yaml:"kms"`
// 	} `yaml:"securityDefinitions"`
// 	Paths struct {
// 		NodeInfo struct {
// 			Get struct {
// 				Description string   `yaml:"description"`
// 				Summary     string   `yaml:"summary"`
// 				Tags        []string `yaml:"tags"`
// 				Produces    []string `yaml:"produces"`
// 				Responses   struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type       string `yaml:"type"`
// 							Properties struct {
// 								ApplicationVersion struct {
// 									Properties struct {
// 										BuildTags struct {
// 											Type string `yaml:"type"`
// 										} `yaml:"build_tags"`
// 										ClientName struct {
// 											Type string `yaml:"type"`
// 										} `yaml:"client_name"`
// 										Commit struct {
// 											Type string `yaml:"type"`
// 										} `yaml:"commit"`
// 										Go struct {
// 											Type string `yaml:"type"`
// 										} `yaml:"go"`
// 										Name struct {
// 											Type string `yaml:"type"`
// 										} `yaml:"name"`
// 										ServerName struct {
// 											Type string `yaml:"type"`
// 										} `yaml:"server_name"`
// 										Version struct {
// 											Type string `yaml:"type"`
// 										} `yaml:"version"`
// 									} `yaml:"properties"`
// 								} `yaml:"application_version"`
// 								NodeInfo struct {
// 									Properties struct {
// 										ID struct {
// 											Type string `yaml:"type"`
// 										} `yaml:"id"`
// 										Moniker struct {
// 											Type    string `yaml:"type"`
// 											Example string `yaml:"example"`
// 										} `yaml:"moniker"`
// 										ProtocolVersion struct {
// 											Properties struct {
// 												P2P struct {
// 													Type    string `yaml:"type"`
// 													Example int    `yaml:"example"`
// 												} `yaml:"p2p"`
// 												Block struct {
// 													Type    string `yaml:"type"`
// 													Example int    `yaml:"example"`
// 												} `yaml:"block"`
// 												App struct {
// 													Type    string `yaml:"type"`
// 													Example int    `yaml:"example"`
// 												} `yaml:"app"`
// 											} `yaml:"properties"`
// 										} `yaml:"protocol_version"`
// 										Network struct {
// 											Type    string `yaml:"type"`
// 											Example string `yaml:"example"`
// 										} `yaml:"network"`
// 										Channels struct {
// 											Type string `yaml:"type"`
// 										} `yaml:"channels"`
// 										ListenAddr struct {
// 											Type    string `yaml:"type"`
// 											Example string `yaml:"example"`
// 										} `yaml:"listen_addr"`
// 										Version struct {
// 											Description string `yaml:"description"`
// 											Type        string `yaml:"type"`
// 											Example     string `yaml:"example"`
// 										} `yaml:"version"`
// 										Other struct {
// 											Description string `yaml:"description"`
// 											Type        string `yaml:"type"`
// 											Properties  struct {
// 												TxIndex struct {
// 													Type    string `yaml:"type"`
// 													Example string `yaml:"example"`
// 												} `yaml:"tx_index"`
// 												RPCAddress struct {
// 													Type    string `yaml:"type"`
// 													Example string `yaml:"example"`
// 												} `yaml:"rpc_address"`
// 											} `yaml:"properties"`
// 										} `yaml:"other"`
// 									} `yaml:"properties"`
// 								} `yaml:"node_info"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/node_info"`
// 		Syncing struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Tags        []string `yaml:"tags"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Responses   struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type       string `yaml:"type"`
// 							Properties struct {
// 								Syncing struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"syncing"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/syncing"`
// 		BlocksLatest struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/blocks/latest"`
// 		BlocksHeight struct {
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 					XExample    int    `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num404 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"404"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/blocks/{height}"`
// 		ValidatorsetsLatest struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type       string `yaml:"type"`
// 							Properties struct {
// 								BlockHeight struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"block_height"`
// 								Validators struct {
// 									Type  string `yaml:"type"`
// 									Items struct {
// 										Ref string `yaml:"$ref"`
// 									} `yaml:"items"`
// 								} `yaml:"validators"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/validatorsets/latest"`
// 		ValidatorsetsHeight struct {
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 					XExample    int    `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type       string `yaml:"type"`
// 							Properties struct {
// 								BlockHeight struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"block_height"`
// 								Validators struct {
// 									Type  string `yaml:"type"`
// 									Items struct {
// 										Ref string `yaml:"$ref"`
// 									} `yaml:"items"`
// 								} `yaml:"validators"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num404 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"404"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/validatorsets/{height}"`
// 		TxsHash struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Tags        []string `yaml:"tags"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Parameters  []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 					XExample    string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/txs/{hash}"`
// 		Txs struct {
// 			Get struct {
// 				Tags        []string `yaml:"tags"`
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Parameters  []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Type        string `yaml:"type"`
// 					Description string `yaml:"description"`
// 					XExample    string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 			Post struct {
// 				Tags        []string `yaml:"tags"`
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Parameters  []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							Tx struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"tx"`
// 							Mode struct {
// 								Type    string `yaml:"type"`
// 								Example string `yaml:"example"`
// 							} `yaml:"mode"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/txs"`
// 		TxsEncode struct {
// 			Post struct {
// 				Tags        []string `yaml:"tags"`
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Parameters  []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							Tx struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"tx"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type       string `yaml:"type"`
// 							Properties struct {
// 								Tx struct {
// 									Type    string `yaml:"type"`
// 									Example string `yaml:"example"`
// 								} `yaml:"tx"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/txs/encode"`
// 		TxsEstimateFee struct {
// 			Post struct {
// 				Tags        []string `yaml:"tags"`
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Parameters  []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							Tx struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"tx"`
// 							GasAdjustment struct {
// 								Type        string `yaml:"type"`
// 								Description string `yaml:"description"`
// 								Example     string `yaml:"example"`
// 							} `yaml:"gas_adjustment"`
// 							GasPrices struct {
// 								Description string `yaml:"description"`
// 								Type        string `yaml:"type"`
// 								Items       struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"items"`
// 							} `yaml:"gas_prices"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/txs/estimate_fee"`
// 		BankBalancesAddress struct {
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 					XExample    string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/bank/balances/{address}"`
// 		BankAccountsAddressTransfers struct {
// 			Post struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Consumes   []string `yaml:"consumes"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type,omitempty"`
// 					XExample    string `yaml:"x-example,omitempty"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 							Coins struct {
// 								Type  string `yaml:"type"`
// 								Items struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"items"`
// 							} `yaml:"coins"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema,omitempty"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num202 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"202"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/bank/accounts/{address}/transfers"`
// 		AuthAccountsAddress struct {
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 					XExample    string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type       string `yaml:"type"`
// 							Properties struct {
// 								Account struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"Account"`
// 								LazyGradedVestingAccount struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"LazyGradedVestingAccount"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num204 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"204"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/auth/accounts/{address}"`
// 		AuthAccountsAddressMultisign struct {
// 			Post struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type,omitempty"`
// 					XExample    string `yaml:"x-example,omitempty"`
// 					Schema      struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"schema,omitempty"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type       string `yaml:"type"`
// 							Properties struct {
// 								MultiSignedTx struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"MultiSignedTx"`
// 								MultiSig struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"MultiSig"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/auth/accounts/{address}/multisign"`
// 		StakingDelegatorsDelegatorAddrDelegations struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 			Post struct {
// 				Summary    string `yaml:"summary"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 							DelegatorAddress struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"delegator_address"`
// 							ValidatorAddress struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"validator_address"`
// 							Amount struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"amount"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Tags      []string `yaml:"tags"`
// 				Consumes  []string `yaml:"consumes"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num401 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"401"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/staking/delegators/{delegatorAddr}/delegations"`
// 		StakingDelegatorsDelegatorAddrDelegationsValidatorAddr struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/staking/delegators/{delegatorAddr}/delegations/{validatorAddr}"`
// 		StakingDelegatorsDelegatorAddrUnbondingDelegations struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 			Post struct {
// 				Summary    string `yaml:"summary"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 							DelegatorAddress struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"delegator_address"`
// 							ValidatorAddress struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"validator_address"`
// 							Amount struct {
// 								Type    string `yaml:"type"`
// 								Example string `yaml:"example"`
// 							} `yaml:"amount"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Tags      []string `yaml:"tags"`
// 				Consumes  []string `yaml:"consumes"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num401 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"401"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/staking/delegators/{delegatorAddr}/unbonding_delegations"`
// 		StakingDelegatorsDelegatorAddrUnbondingDelegationsValidatorAddr struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/staking/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr}"`
// 		StakingRedelegations struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/staking/redelegations"`
// 		StakingDelegatorsDelegatorAddrValidators struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/staking/delegators/{delegatorAddr}/validators"`
// 		StakingDelegatorsDelegatorAddrValidatorsValidatorAddr struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/staking/delegators/{delegatorAddr}/validators/{validatorAddr}"`
// 		StakingDelegatorsDelegatorAddrTxs struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num204 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"204"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/staking/delegators/{delegatorAddr}/txs"`
// 		StakingValidators struct {
// 			Get struct {
// 				Summary    string `yaml:"summary"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Type        string `yaml:"type"`
// 					Description string `yaml:"description"`
// 					XExample    string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/staking/validators"`
// 		StakingValidatorsValidatorAddr struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/staking/validators/{validatorAddr}"`
// 		StakingValidatorsValidatorAddrDelegations struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/staking/validators/{validatorAddr}/delegations"`
// 		StakingValidatorsValidatorAddrUnbondingDelegations struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/staking/validators/{validatorAddr}/unbonding_delegations"`
// 		StakingPool struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type       string `yaml:"type"`
// 							Properties struct {
// 								BondedTokens struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"bonded_tokens"`
// 								NotBondedTokens struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"not_bonded_tokens"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/staking/pool"`
// 		StakingParameters struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type       string `yaml:"type"`
// 							Properties struct {
// 								UnbondingTime struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"unbonding_time"`
// 								MaxValidators struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"max_validators"`
// 								MaxEntries struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"max_entries"`
// 								BondDenom struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"bond_denom"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/staking/parameters"`
// 		SlashingValidatorsValidatorPubKeySigningInfo struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Type        string `yaml:"type"`
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					Required    bool   `yaml:"required"`
// 					In          string `yaml:"in"`
// 					XExample    string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num204 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"204"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/slashing/validators/{validatorPubKey}/signing_info"`
// 		SlashingSigningInfos struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Type        string `yaml:"type"`
// 					Required    bool   `yaml:"required"`
// 					XExample    int    `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/slashing/signing_infos"`
// 		SlashingValidatorsValidatorAddrUnjail struct {
// 			Post struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Type        string `yaml:"type,omitempty"`
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					Required    bool   `yaml:"required"`
// 					In          string `yaml:"in"`
// 					XExample    string `yaml:"x-example,omitempty"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema,omitempty"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/slashing/validators/{validatorAddr}/unjail"`
// 		SlashingParameters struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type       string `yaml:"type"`
// 							Properties struct {
// 								MaxEvidenceAge struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"max_evidence_age"`
// 								SignedBlocksWindow struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"signed_blocks_window"`
// 								MinSignedPerWindow struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"min_signed_per_window"`
// 								DowntimeJailDuration struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"downtime_jail_duration"`
// 								SlashFractionDoubleSign struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"slash_fraction_double_sign"`
// 								SlashFractionDowntime struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"slash_fraction_downtime"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/slashing/parameters"`
// 		GovProposals struct {
// 			Post struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					In          string `yaml:"in"`
// 					Required    bool   `yaml:"required"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 							Title struct {
// 								Type string `yaml:"type"`
// 							} `yaml:"title"`
// 							Description struct {
// 								Type string `yaml:"type"`
// 							} `yaml:"description"`
// 							ProposalType struct {
// 								Type    string `yaml:"type"`
// 								Example string `yaml:"example"`
// 							} `yaml:"proposal_type"`
// 							Proposer struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"proposer"`
// 							InitialDeposit struct {
// 								Type  string `yaml:"type"`
// 								Items struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"items"`
// 							} `yaml:"initial_deposit"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/gov/proposals"`
// 		GovProposalsParamChange struct {
// 			Post struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					In          string `yaml:"in"`
// 					Required    bool   `yaml:"required"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 							Title struct {
// 								Type     string `yaml:"type"`
// 								XExample string `yaml:"x-example"`
// 							} `yaml:"title"`
// 							Description struct {
// 								Type     string `yaml:"type"`
// 								XExample string `yaml:"x-example"`
// 							} `yaml:"description"`
// 							Proposer struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"proposer"`
// 							Deposit struct {
// 								Type  string `yaml:"type"`
// 								Items struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"items"`
// 							} `yaml:"deposit"`
// 							Changes struct {
// 								Type  string `yaml:"type"`
// 								Items struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"items"`
// 							} `yaml:"changes"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/gov/proposals/param_change"`
// 		GovProposalsCommunityPoolSpend struct {
// 			Post struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					In          string `yaml:"in"`
// 					Required    bool   `yaml:"required"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 							Title struct {
// 								Type     string `yaml:"type"`
// 								XExample string `yaml:"x-example"`
// 							} `yaml:"title"`
// 							Description struct {
// 								Type     string `yaml:"type"`
// 								XExample string `yaml:"x-example"`
// 							} `yaml:"description"`
// 							Proposer struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"proposer"`
// 							Deposit struct {
// 								Type  string `yaml:"type"`
// 								Items struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"items"`
// 							} `yaml:"deposit"`
// 							Recipient struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"recipient"`
// 							Amount struct {
// 								Type  string `yaml:"type"`
// 								Items struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"items"`
// 							} `yaml:"amount"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/gov/proposals/community_pool_spend"`
// 		GovProposalsTaxRateUpdate struct {
// 			Post struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					In          string `yaml:"in"`
// 					Required    bool   `yaml:"required"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 							Title struct {
// 								Type     string `yaml:"type"`
// 								XExample string `yaml:"x-example"`
// 							} `yaml:"title"`
// 							Description struct {
// 								Type     string `yaml:"type"`
// 								XExample string `yaml:"x-example"`
// 							} `yaml:"description"`
// 							Proposer struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"proposer"`
// 							Deposit struct {
// 								Type  string `yaml:"type"`
// 								Items struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"items"`
// 							} `yaml:"deposit"`
// 							TaxRate struct {
// 								Type    string `yaml:"type"`
// 								Format  string `yaml:"format"`
// 								Example string `yaml:"example"`
// 							} `yaml:"tax_rate"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/gov/proposals/tax_rate_update"`
// 		GovProposalsRewardWeightUpdate struct {
// 			Post struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					In          string `yaml:"in"`
// 					Required    bool   `yaml:"required"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 							Title struct {
// 								Type     string `yaml:"type"`
// 								XExample string `yaml:"x-example"`
// 							} `yaml:"title"`
// 							Description struct {
// 								Type     string `yaml:"type"`
// 								XExample string `yaml:"x-example"`
// 							} `yaml:"description"`
// 							Proposer struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"proposer"`
// 							Deposit struct {
// 								Type  string `yaml:"type"`
// 								Items struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"items"`
// 							} `yaml:"deposit"`
// 							RewardWeight struct {
// 								Type    string `yaml:"type"`
// 								Format  string `yaml:"format"`
// 								Example string `yaml:"example"`
// 							} `yaml:"reward_weight"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/gov/proposals/reward_weight_update"`
// 		GovProposalsProposalID struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Type     string `yaml:"type"`
// 					Name     string `yaml:"name"`
// 					Required bool   `yaml:"required"`
// 					In       string `yaml:"in"`
// 					XExample string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/gov/proposals/{proposalId}"`
// 		GovProposalsProposalIDProposer struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Type     string `yaml:"type"`
// 					Name     string `yaml:"name"`
// 					Required bool   `yaml:"required"`
// 					In       string `yaml:"in"`
// 					XExample string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/gov/proposals/{proposalId}/proposer"`
// 		GovProposalsProposalIDDeposits struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Type     string `yaml:"type"`
// 					Name     string `yaml:"name"`
// 					Required bool   `yaml:"required"`
// 					In       string `yaml:"in"`
// 					XExample string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 			Post struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Type        string `yaml:"type,omitempty"`
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					Required    bool   `yaml:"required"`
// 					In          string `yaml:"in"`
// 					XExample    string `yaml:"x-example,omitempty"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 							Depositor struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"depositor"`
// 							Amount struct {
// 								Type  string `yaml:"type"`
// 								Items struct {
// 									Ref string `yaml:"$ref"`
// 								} `yaml:"items"`
// 							} `yaml:"amount"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema,omitempty"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num401 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"401"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/gov/proposals/{proposalId}/deposits"`
// 		GovProposalsProposalIDDepositsDepositor struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Type        string `yaml:"type"`
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					Required    bool   `yaml:"required"`
// 					In          string `yaml:"in"`
// 					XExample    string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num404 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"404"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/gov/proposals/{proposalId}/deposits/{depositor}"`
// 		GovProposalsProposalIDVotes struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Type        string `yaml:"type"`
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					Required    bool   `yaml:"required"`
// 					In          string `yaml:"in"`
// 					XExample    string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 			Post struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Type        string `yaml:"type,omitempty"`
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					Required    bool   `yaml:"required"`
// 					In          string `yaml:"in"`
// 					XExample    string `yaml:"x-example,omitempty"`
// 					Schema      struct {
// 						Type       string `yaml:"type"`
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 							Voter struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"voter"`
// 							Option struct {
// 								Type    string `yaml:"type"`
// 								Example string `yaml:"example"`
// 							} `yaml:"option"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema,omitempty"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num401 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"401"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/gov/proposals/{proposalId}/votes"`
// 		GovProposalsProposalIDVotesVoter struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Type        string `yaml:"type"`
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					Required    bool   `yaml:"required"`
// 					In          string `yaml:"in"`
// 					XExample    string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num404 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"404"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/gov/proposals/{proposalId}/votes/{voter}"`
// 		GovProposalsProposalIDTally struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Parameters  []struct {
// 					Type        string `yaml:"type"`
// 					Description string `yaml:"description"`
// 					Name        string `yaml:"name"`
// 					Required    bool   `yaml:"required"`
// 					In          string `yaml:"in"`
// 					XExample    string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/gov/proposals/{proposalId}/tally"`
// 		GovParametersDeposit struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Responses   struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type       string `yaml:"type"`
// 							Properties struct {
// 								MinDeposit struct {
// 									Type  string `yaml:"type"`
// 									Items struct {
// 										Ref string `yaml:"$ref"`
// 									} `yaml:"items"`
// 								} `yaml:"min_deposit"`
// 								MaxDepositPeriod struct {
// 									Type    string `yaml:"type"`
// 									Example string `yaml:"example"`
// 								} `yaml:"max_deposit_period"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num404 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"404"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/gov/parameters/deposit"`
// 		GovParametersTallying struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Responses   struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Properties struct {
// 								Quorum struct {
// 									Type    string `yaml:"type"`
// 									Example string `yaml:"example"`
// 								} `yaml:"quorum"`
// 								Threshold struct {
// 									Type    string `yaml:"type"`
// 									Example string `yaml:"example"`
// 								} `yaml:"threshold"`
// 								Veto struct {
// 									Type    string `yaml:"type"`
// 									Example string `yaml:"example"`
// 								} `yaml:"veto"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num404 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"404"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/gov/parameters/tallying"`
// 		GovParametersVoting struct {
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Responses   struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Properties struct {
// 								VotingPeriod struct {
// 									Type    string `yaml:"type"`
// 									Example string `yaml:"example"`
// 								} `yaml:"voting_period"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num404 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"404"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/gov/parameters/voting"`
// 		DistributionDelegatorsDelegatorAddrRewards struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Produces    []string `yaml:"produces"`
// 				Tags        []string `yaml:"tags"`
// 				Responses   struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 			Post struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Tags        []string `yaml:"tags"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Parameters  []struct {
// 					In     string `yaml:"in"`
// 					Name   string `yaml:"name"`
// 					Schema struct {
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num401 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"401"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/distribution/delegators/{delegatorAddr}/rewards"`
// 		DistributionDelegatorsDelegatorAddrRewardsValidatorAddr struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Tags        []string `yaml:"tags"`
// 				Produces    []string `yaml:"produces"`
// 				Responses   struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 			Post struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Tags        []string `yaml:"tags"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Parameters  []struct {
// 					In     string `yaml:"in"`
// 					Name   string `yaml:"name"`
// 					Schema struct {
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num401 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"401"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/distribution/delegators/{delegatorAddr}/rewards/{validatorAddr}"`
// 		DistributionDelegatorsDelegatorAddrWithdrawAddress struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Tags        []string `yaml:"tags"`
// 				Produces    []string `yaml:"produces"`
// 				Responses   struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 			Post struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Tags        []string `yaml:"tags"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Parameters  []struct {
// 					In     string `yaml:"in"`
// 					Name   string `yaml:"name"`
// 					Schema struct {
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 							WithdrawAddress struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"withdraw_address"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num401 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"401"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/distribution/delegators/{delegatorAddr}/withdraw_address"`
// 		DistributionValidatorsValidatorAddr struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Tags        []string `yaml:"tags"`
// 				Produces    []string `yaml:"produces"`
// 				Responses   struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/distribution/validators/{validatorAddr}"`
// 		DistributionValidatorsValidatorAddrOutstandingRewards struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/distribution/validators/{validatorAddr}/outstanding_rewards"`
// 		DistributionValidatorsValidatorAddrRewards struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Tags        []string `yaml:"tags"`
// 				Produces    []string `yaml:"produces"`
// 				Responses   struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 			Post struct {
// 				Summary     string   `yaml:"summary"`
// 				Description string   `yaml:"description"`
// 				Tags        []string `yaml:"tags"`
// 				Consumes    []string `yaml:"consumes"`
// 				Produces    []string `yaml:"produces"`
// 				Parameters  []struct {
// 					In     string `yaml:"in"`
// 					Name   string `yaml:"name"`
// 					Schema struct {
// 						Properties struct {
// 							BaseReq struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"base_req"`
// 						} `yaml:"properties"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num401 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"401"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 		} `yaml:"/distribution/validators/{validatorAddr}/rewards"`
// 		DistributionCommunityPool struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/distribution/community_pool"`
// 		DistributionParams struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Properties struct {
// 								BaseProposerReward struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"base_proposer_reward"`
// 								BonusProposerReward struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"bonus_proposer_reward"`
// 								CommunityTax struct {
// 									Type string `yaml:"type"`
// 								} `yaml:"community_tax"`
// 							} `yaml:"properties"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/distribution/params"`
// 		SupplyTotal struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/supply/total"`
// 		SupplyTotalDenomination struct {
// 			Parameters []struct {
// 				In          string `yaml:"in"`
// 				Name        string `yaml:"name"`
// 				Description string `yaml:"description"`
// 				Required    bool   `yaml:"required"`
// 				Type        string `yaml:"type"`
// 				XExample    string `yaml:"x-example"`
// 			} `yaml:"parameters"`
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type string `yaml:"type"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/supply/total/{denomination}"`
// 		MarketSwap struct {
// 			Post struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In     string `yaml:"in"`
// 					Name   string `yaml:"name"`
// 					Schema struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"schema"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Type        string `yaml:"type"`
// 					Required    bool   `yaml:"required"`
// 					XExample    string `yaml:"x-example"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/market/swap"`
// 		MarketTerraPoolDelta struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type    string `yaml:"type"`
// 							Format  string `yaml:"format"`
// 							Example string `yaml:"example"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/market/terra_pool_delta"`
// 		MarketParameters struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/market/parameters"`
// 		OracleDenomsDenomVotes struct {
// 			Post struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description,omitempty"`
// 					Required    bool   `yaml:"required,omitempty"`
// 					Type        string `yaml:"type,omitempty"`
// 					Schema      struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"schema,omitempty"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/oracle/denoms/{denom}/votes"`
// 		OracleDenomsDenomVotesValidator struct {
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description,omitempty"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/oracle/denoms/{denom}/votes/{validator}"`
// 		OracleVotersValidatorVotes struct {
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In       string `yaml:"in"`
// 					Name     string `yaml:"name"`
// 					Required bool   `yaml:"required"`
// 					Type     string `yaml:"type"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/oracle/voters/{validator}/votes"`
// 		OracleDenomsDenomPrevotes struct {
// 			Post struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description,omitempty"`
// 					Required    bool   `yaml:"required,omitempty"`
// 					Type        string `yaml:"type,omitempty"`
// 					Schema      struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"schema,omitempty"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/oracle/denoms/{denom}/prevotes"`
// 		OracleDenomsDenomPrevotesValidator struct {
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description,omitempty"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/oracle/denoms/{denom}/prevotes/{validator}"`
// 		OracleVotersValidatorPrevotes struct {
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In       string `yaml:"in"`
// 					Name     string `yaml:"name"`
// 					Required bool   `yaml:"required"`
// 					Type     string `yaml:"type"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/oracle/voters/{validator}/prevotes"`
// 		OracleDenomsDenomExchangeRate struct {
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/oracle/denoms/{denom}/exchange_rate"`
// 		OracleDenomsExchangeRates struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/oracle/denoms/exchange_rates"`
// 		OracleDenomsActives struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Type string `yaml:"type"`
// 							} `yaml:"items"`
// 							Example []string `yaml:"example"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/oracle/denoms/actives"`
// 		OracleVotersValidatorFeeder struct {
// 			Post struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description,omitempty"`
// 					Required    bool   `yaml:"required,omitempty"`
// 					Type        string `yaml:"type,omitempty"`
// 					Schema      struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"schema,omitempty"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"post"`
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/oracle/voters/{validator}/feeder"`
// 		OracleVotersValidatorMiss struct {
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type    string `yaml:"type"`
// 							Format  string `yaml:"format"`
// 							Example string `yaml:"example"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/oracle/voters/{validator}/miss"`
// 		OracleParameters struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num400 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"400"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/oracle/parameters"`
// 		TreasuryTaxRate struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type    string `yaml:"type"`
// 							Format  string `yaml:"format"`
// 							Example string `yaml:"example"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/treasury/tax_rate"`
// 		TreasuryTaxCapDenom struct {
// 			Get struct {
// 				Summary    string   `yaml:"summary"`
// 				Tags       []string `yaml:"tags"`
// 				Produces   []string `yaml:"produces"`
// 				Parameters []struct {
// 					In          string `yaml:"in"`
// 					Name        string `yaml:"name"`
// 					Description string `yaml:"description"`
// 					Required    bool   `yaml:"required"`
// 					Type        string `yaml:"type"`
// 				} `yaml:"parameters"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num404 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"404"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/treasury/tax_cap/{denom}"`
// 		TreasuryRewardWeight struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type        string `yaml:"type"`
// 							Format      string `yaml:"format"`
// 							Example     string `yaml:"example"`
// 							Description string `yaml:"description"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/treasury/reward_weight"`
// 		TreasuryTaxProceeds struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/treasury/tax_proceeds"`
// 		TreasurySeigniorageProceeds struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num500 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"500"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/treasury/seigniorage_proceeds"`
// 		TreasuryParameters struct {
// 			Get struct {
// 				Summary   string   `yaml:"summary"`
// 				Tags      []string `yaml:"tags"`
// 				Produces  []string `yaml:"produces"`
// 				Responses struct {
// 					Num200 struct {
// 						Description string `yaml:"description"`
// 						Schema      struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"schema"`
// 					} `yaml:"200"`
// 					Num404 struct {
// 						Description string `yaml:"description"`
// 					} `yaml:"404"`
// 				} `yaml:"responses"`
// 			} `yaml:"get"`
// 		} `yaml:"/treasury/parameters"`
// 	} `yaml:"paths"`
// 	Definitions struct {
// 		CheckTxResult struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Code struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"code"`
// 				Data struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"data"`
// 				GasUsed struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"gas_used"`
// 				GasWanted struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"gas_wanted"`
// 				Info struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"info"`
// 				Log struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"log"`
// 				Tags struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"tags"`
// 			} `yaml:"properties"`
// 			Example struct {
// 				Code      int      `yaml:"code"`
// 				Data      string   `yaml:"data"`
// 				Log       string   `yaml:"log"`
// 				GasUsed   int      `yaml:"gas_used"`
// 				GasWanted int      `yaml:"gas_wanted"`
// 				Info      string   `yaml:"info"`
// 				Tags      []string `yaml:"tags"`
// 			} `yaml:"example"`
// 		} `yaml:"CheckTxResult"`
// 		DeliverTxResult struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Code struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"code"`
// 				Data struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"data"`
// 				GasUsed struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"gas_used"`
// 				GasWanted struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"gas_wanted"`
// 				Info struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"info"`
// 				Log struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"log"`
// 				Tags struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"tags"`
// 			} `yaml:"properties"`
// 			Example struct {
// 				Code      int      `yaml:"code"`
// 				Data      string   `yaml:"data"`
// 				Log       string   `yaml:"log"`
// 				GasUsed   int      `yaml:"gas_used"`
// 				GasWanted int      `yaml:"gas_wanted"`
// 				Info      string   `yaml:"info"`
// 				Tags      []string `yaml:"tags"`
// 			} `yaml:"example"`
// 		} `yaml:"DeliverTxResult"`
// 		BroadcastTxCommitResult struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				CheckTx struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"check_tx"`
// 				DeliverTx struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"deliver_tx"`
// 				Hash struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"hash"`
// 				Height struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"height"`
// 			} `yaml:"properties"`
// 		} `yaml:"BroadcastTxCommitResult"`
// 		KVPair struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Key struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"key"`
// 				Value struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"value"`
// 			} `yaml:"properties"`
// 		} `yaml:"KVPair"`
// 		Msg struct {
// 			Type string `yaml:"type"`
// 		} `yaml:"Msg"`
// 		Address struct {
// 			Type        string `yaml:"type"`
// 			Description string `yaml:"description"`
// 			Example     string `yaml:"example"`
// 		} `yaml:"Address"`
// 		ValidatorAddress struct {
// 			Type        string `yaml:"type"`
// 			Description string `yaml:"description"`
// 			Example     string `yaml:"example"`
// 		} `yaml:"ValidatorAddress"`
// 		Coin struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Denom struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"denom"`
// 				Amount struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"amount"`
// 			} `yaml:"properties"`
// 		} `yaml:"Coin"`
// 		DecCoin struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Denom struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"denom"`
// 				Amount struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"amount"`
// 			} `yaml:"properties"`
// 		} `yaml:"DecCoin"`
// 		Hash struct {
// 			Type    string `yaml:"type"`
// 			Example string `yaml:"example"`
// 		} `yaml:"Hash"`
// 		TxQuery struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Hash struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"hash"`
// 				Height struct {
// 					Type    string `yaml:"type"`
// 					Example int    `yaml:"example"`
// 				} `yaml:"height"`
// 				Tx struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"tx"`
// 				Result struct {
// 					Type       string `yaml:"type"`
// 					Properties struct {
// 						Log struct {
// 							Type string `yaml:"type"`
// 						} `yaml:"log"`
// 						GasWanted struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"gas_wanted"`
// 						GasUsed struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"gas_used"`
// 						Tags struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"tags"`
// 					} `yaml:"properties"`
// 				} `yaml:"result"`
// 			} `yaml:"properties"`
// 		} `yaml:"TxQuery"`
// 		PaginatedQueryTxs struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				TotalCount struct {
// 					Type    string `yaml:"type"`
// 					Example int    `yaml:"example"`
// 				} `yaml:"total_count"`
// 				Count struct {
// 					Type    string `yaml:"type"`
// 					Example int    `yaml:"example"`
// 				} `yaml:"count"`
// 				PageNumber struct {
// 					Type    string `yaml:"type"`
// 					Example int    `yaml:"example"`
// 				} `yaml:"page_number"`
// 				PageTotal struct {
// 					Type    string `yaml:"type"`
// 					Example int    `yaml:"example"`
// 				} `yaml:"page_total"`
// 				Limit struct {
// 					Type    string `yaml:"type"`
// 					Example int    `yaml:"example"`
// 				} `yaml:"limit"`
// 				Txs struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"txs"`
// 			} `yaml:"properties"`
// 		} `yaml:"PaginatedQueryTxs"`
// 		StdTx struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Msg struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"msg"`
// 				Fee struct {
// 					Type       string `yaml:"type"`
// 					Properties struct {
// 						Gas struct {
// 							Type string `yaml:"type"`
// 						} `yaml:"gas"`
// 						Amount struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"amount"`
// 					} `yaml:"properties"`
// 				} `yaml:"fee"`
// 				Memo struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"memo"`
// 				Signature struct {
// 					Type       string `yaml:"type"`
// 					Properties struct {
// 						Signature struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"signature"`
// 						PubKey struct {
// 							Type       string `yaml:"type"`
// 							Properties struct {
// 								Type struct {
// 									Type    string `yaml:"type"`
// 									Example string `yaml:"example"`
// 								} `yaml:"type"`
// 								Value struct {
// 									Type    string `yaml:"type"`
// 									Example string `yaml:"example"`
// 								} `yaml:"value"`
// 							} `yaml:"properties"`
// 						} `yaml:"pub_key"`
// 						AccountNumber struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"account_number"`
// 						Sequence struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"sequence"`
// 					} `yaml:"properties"`
// 				} `yaml:"signature"`
// 			} `yaml:"properties"`
// 		} `yaml:"StdTx"`
// 		UnsignedStdTx struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Msg struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"msg"`
// 				Fee struct {
// 					Type       string `yaml:"type"`
// 					Properties struct {
// 						Gas struct {
// 							Type string `yaml:"type"`
// 						} `yaml:"gas"`
// 						Amount struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Ref string `yaml:"$ref"`
// 							} `yaml:"items"`
// 						} `yaml:"amount"`
// 					} `yaml:"properties"`
// 				} `yaml:"fee"`
// 				Memo struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"memo"`
// 				Signatures struct {
// 					Type    string      `yaml:"type"`
// 					Example interface{} `yaml:"example"`
// 				} `yaml:"signatures"`
// 			} `yaml:"properties"`
// 		} `yaml:"UnsignedStdTx"`
// 		BlockID struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Hash struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"hash"`
// 				Parts struct {
// 					Type       string `yaml:"type"`
// 					Properties struct {
// 						Total struct {
// 							Type    string `yaml:"type"`
// 							Example int    `yaml:"example"`
// 						} `yaml:"total"`
// 						Hash struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"hash"`
// 					} `yaml:"properties"`
// 				} `yaml:"parts"`
// 			} `yaml:"properties"`
// 		} `yaml:"BlockID"`
// 		BlockHeader struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				ChainID struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"chain_id"`
// 				Height struct {
// 					Type    string `yaml:"type"`
// 					Example int    `yaml:"example"`
// 				} `yaml:"height"`
// 				Time struct {
// 					Type    string    `yaml:"type"`
// 					Example time.Time `yaml:"example"`
// 				} `yaml:"time"`
// 				NumTxs struct {
// 					Type    string `yaml:"type"`
// 					Example int    `yaml:"example"`
// 				} `yaml:"num_txs"`
// 				LastBlockID struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"last_block_id"`
// 				TotalTxs struct {
// 					Type    string `yaml:"type"`
// 					Example int    `yaml:"example"`
// 				} `yaml:"total_txs"`
// 				LastCommitHash struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"last_commit_hash"`
// 				DataHash struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"data_hash"`
// 				ValidatorsHash struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"validators_hash"`
// 				NextValidatorsHash struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"next_validators_hash"`
// 				ConsensusHash struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"consensus_hash"`
// 				AppHash struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"app_hash"`
// 				LastResultsHash struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"last_results_hash"`
// 				EvidenceHash struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"evidence_hash"`
// 				ProposerAddress struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"proposer_address"`
// 				Version struct {
// 					Type       string `yaml:"type"`
// 					Properties struct {
// 						Block struct {
// 							Type    string `yaml:"type"`
// 							Example int    `yaml:"example"`
// 						} `yaml:"block"`
// 						App struct {
// 							Type    string `yaml:"type"`
// 							Example int    `yaml:"example"`
// 						} `yaml:"app"`
// 					} `yaml:"properties"`
// 				} `yaml:"version"`
// 			} `yaml:"properties"`
// 		} `yaml:"BlockHeader"`
// 		Block struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Header struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"header"`
// 				Txs struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Type string `yaml:"type"`
// 					} `yaml:"items"`
// 				} `yaml:"txs"`
// 				Evidence struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Type string `yaml:"type"`
// 					} `yaml:"items"`
// 				} `yaml:"evidence"`
// 				LastCommit struct {
// 					Type       string `yaml:"type"`
// 					Properties struct {
// 						BlockID struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"block_id"`
// 						Precommits struct {
// 							Type  string `yaml:"type"`
// 							Items struct {
// 								Type       string `yaml:"type"`
// 								Properties struct {
// 									ValidatorAddress struct {
// 										Type string `yaml:"type"`
// 									} `yaml:"validator_address"`
// 									ValidatorIndex struct {
// 										Type    string `yaml:"type"`
// 										Example string `yaml:"example"`
// 									} `yaml:"validator_index"`
// 									Height struct {
// 										Type    string `yaml:"type"`
// 										Example string `yaml:"example"`
// 									} `yaml:"height"`
// 									Round struct {
// 										Type    string `yaml:"type"`
// 										Example string `yaml:"example"`
// 									} `yaml:"round"`
// 									Timestamp struct {
// 										Type    string    `yaml:"type"`
// 										Example time.Time `yaml:"example"`
// 									} `yaml:"timestamp"`
// 									Type struct {
// 										Type    string `yaml:"type"`
// 										Example int    `yaml:"example"`
// 									} `yaml:"type"`
// 									BlockID struct {
// 										Ref string `yaml:"$ref"`
// 									} `yaml:"block_id"`
// 									Signature struct {
// 										Type    string `yaml:"type"`
// 										Example string `yaml:"example"`
// 									} `yaml:"signature"`
// 								} `yaml:"properties"`
// 							} `yaml:"items"`
// 						} `yaml:"precommits"`
// 					} `yaml:"properties"`
// 				} `yaml:"last_commit"`
// 			} `yaml:"properties"`
// 		} `yaml:"Block"`
// 		BlockQuery struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				BlockMeta struct {
// 					Type       string `yaml:"type"`
// 					Properties struct {
// 						Header struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"header"`
// 						BlockID struct {
// 							Ref string `yaml:"$ref"`
// 						} `yaml:"block_id"`
// 					} `yaml:"properties"`
// 				} `yaml:"block_meta"`
// 				Block struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"block"`
// 			} `yaml:"properties"`
// 		} `yaml:"BlockQuery"`
// 		DelegationDelegatorReward struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				ValidatorAddress struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"validator_address"`
// 				Reward struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"reward"`
// 			} `yaml:"properties"`
// 		} `yaml:"DelegationDelegatorReward"`
// 		DelegatorTotalRewards struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Rewards struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"rewards"`
// 				Total struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"total"`
// 			} `yaml:"properties"`
// 		} `yaml:"DelegatorTotalRewards"`
// 		BaseReq struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				From struct {
// 					Type        string `yaml:"type"`
// 					Example     string `yaml:"example"`
// 					Description string `yaml:"description"`
// 				} `yaml:"from"`
// 				Memo struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"memo"`
// 				ChainID struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"chain_id"`
// 				AccountNumber struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"account_number"`
// 				Sequence struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"sequence"`
// 				Gas struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"gas"`
// 				GasAdjustment struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"gas_adjustment"`
// 				Fees struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"fees"`
// 				Simulate struct {
// 					Type        string `yaml:"type"`
// 					Example     bool   `yaml:"example"`
// 					Description string `yaml:"description"`
// 				} `yaml:"simulate"`
// 			} `yaml:"properties"`
// 		} `yaml:"BaseReq"`
// 		TendermintValidator struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Address struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"address"`
// 				PubKey struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"pub_key"`
// 				VotingPower struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"voting_power"`
// 				ProposerPriority struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"proposer_priority"`
// 			} `yaml:"properties"`
// 		} `yaml:"TendermintValidator"`
// 		TextProposal struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				ProposalID struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"proposal_id"`
// 				Title struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"title"`
// 				Description struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"description"`
// 				ProposalType struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"proposal_type"`
// 				ProposalStatus struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"proposal_status"`
// 				FinalTallyResult struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"final_tally_result"`
// 				SubmitTime struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"submit_time"`
// 				TotalDeposit struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"total_deposit"`
// 				VotingStartTime struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"voting_start_time"`
// 			} `yaml:"properties"`
// 		} `yaml:"TextProposal"`
// 		Proposer struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				ProposalID struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"proposal_id"`
// 				Proposer struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"proposer"`
// 			} `yaml:"properties"`
// 		} `yaml:"Proposer"`
// 		Deposit struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Amount struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"amount"`
// 				ProposalID struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"proposal_id"`
// 				Depositor struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"depositor"`
// 			} `yaml:"properties"`
// 		} `yaml:"Deposit"`
// 		TallyResult struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Yes struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"yes"`
// 				Abstain struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"abstain"`
// 				No struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"no"`
// 				NoWithVeto struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"no_with_veto"`
// 			} `yaml:"properties"`
// 		} `yaml:"TallyResult"`
// 		Vote struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Voter struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"voter"`
// 				ProposalID struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"proposal_id"`
// 				Option struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"option"`
// 			} `yaml:"properties"`
// 		} `yaml:"Vote"`
// 		Validator struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				OperatorAddress struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"operator_address"`
// 				ConsensusPubkey struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"consensus_pubkey"`
// 				Jailed struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"jailed"`
// 				Status struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"status"`
// 				Tokens struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"tokens"`
// 				DelegatorShares struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"delegator_shares"`
// 				Description struct {
// 					Type       string `yaml:"type"`
// 					Properties struct {
// 						Moniker struct {
// 							Type string `yaml:"type"`
// 						} `yaml:"moniker"`
// 						Identity struct {
// 							Type string `yaml:"type"`
// 						} `yaml:"identity"`
// 						Website struct {
// 							Type string `yaml:"type"`
// 						} `yaml:"website"`
// 						Details struct {
// 							Type string `yaml:"type"`
// 						} `yaml:"details"`
// 					} `yaml:"properties"`
// 				} `yaml:"description"`
// 				BondHeight struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"bond_height"`
// 				BondIntraTxCounter struct {
// 					Type    string `yaml:"type"`
// 					Example int    `yaml:"example"`
// 				} `yaml:"bond_intra_tx_counter"`
// 				UnbondingHeight struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"unbonding_height"`
// 				UnbondingTime struct {
// 					Type    string    `yaml:"type"`
// 					Example time.Time `yaml:"example"`
// 				} `yaml:"unbonding_time"`
// 				Commission struct {
// 					Type       string `yaml:"type"`
// 					Properties struct {
// 						Rate struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"rate"`
// 						MaxRate struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"max_rate"`
// 						MaxChangeRate struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"max_change_rate"`
// 						UpdateTime struct {
// 							Type    string    `yaml:"type"`
// 							Example time.Time `yaml:"example"`
// 						} `yaml:"update_time"`
// 					} `yaml:"properties"`
// 				} `yaml:"commission"`
// 			} `yaml:"properties"`
// 		} `yaml:"Validator"`
// 		Delegation struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				DelegatorAddress struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"delegator_address"`
// 				ValidatorAddress struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"validator_address"`
// 				Shares struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"shares"`
// 				Balance struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"balance"`
// 			} `yaml:"properties"`
// 		} `yaml:"Delegation"`
// 		UnbondingDelegation struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				DelegatorAddress struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"delegator_address"`
// 				ValidatorAddress struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"validator_address"`
// 				Entries struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"entries"`
// 			} `yaml:"properties"`
// 		} `yaml:"UnbondingDelegation"`
// 		UnbondingEntry struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				InitialBalance struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"initial_balance"`
// 				Balance struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"balance"`
// 				CreationHeight struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"creation_height"`
// 				CompletionTime struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"completion_time"`
// 			} `yaml:"properties"`
// 		} `yaml:"UnbondingEntry"`
// 		Redelegation struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				DelegatorAddress struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"delegator_address"`
// 				ValidatorSrcAddress struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"validator_src_address"`
// 				ValidatorDstAddress struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"validator_dst_address"`
// 				Entries struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"entries"`
// 			} `yaml:"properties"`
// 		} `yaml:"Redelegation"`
// 		RedelegationEntry struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				CreationHeight struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"creation_height"`
// 				CompletionTime struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"completion_time"`
// 				InitialBalance struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"initial_balance"`
// 				Balance struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"balance"`
// 				SharesDst struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"shares_dst"`
// 			} `yaml:"properties"`
// 		} `yaml:"RedelegationEntry"`
// 		ValidatorDistInfo struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				OperatorAddress struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"operator_address"`
// 				SelfBondRewards struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"self_bond_rewards"`
// 				ValCommission struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"val_commission"`
// 			} `yaml:"properties"`
// 		} `yaml:"ValidatorDistInfo"`
// 		PublicKey struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Type struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"type"`
// 				Value struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"value"`
// 			} `yaml:"properties"`
// 		} `yaml:"PublicKey"`
// 		SigningInfo struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Address struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"address"`
// 				StartHeight struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"start_height"`
// 				IndexOffset struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"index_offset"`
// 				JailedUntil struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"jailed_until"`
// 				Tombstoned struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"tombstoned"`
// 				MissedBlocksCounter struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"missed_blocks_counter"`
// 			} `yaml:"properties"`
// 		} `yaml:"SigningInfo"`
// 		ParamChange struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Subspace struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"subspace"`
// 				Key struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"key"`
// 				Subkey struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"subkey"`
// 				Value struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"value"`
// 			} `yaml:"properties"`
// 		} `yaml:"ParamChange"`
// 		Supply struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Total struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"total"`
// 			} `yaml:"properties"`
// 		} `yaml:"Supply"`
// 		BaseAccount struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				AccountNumber struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"account_number"`
// 				Address struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"address"`
// 				Coins struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"coins"`
// 				PublicKey struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"public_key"`
// 				Sequence struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"sequence"`
// 			} `yaml:"properties"`
// 		} `yaml:"BaseAccount"`
// 		Account struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Type struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"type"`
// 				Value struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"value"`
// 			} `yaml:"properties"`
// 		} `yaml:"Account"`
// 		BaseVestingAccount struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				BaseAccount struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"BaseAccount"`
// 				OriginalVesting struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"original_vesting"`
// 				DelegatedFree struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"delegated_free"`
// 				DelegatedVesting struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"delegated_vesting"`
// 				EndTime struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"end_time"`
// 			} `yaml:"properties"`
// 		} `yaml:"BaseVestingAccount"`
// 		Schedule struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				StartTime struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"start_time"`
// 				EndTime struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"end_time"`
// 				Ratio struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"ratio"`
// 			} `yaml:"properties"`
// 		} `yaml:"Schedule"`
// 		VestingSchedule struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Denom struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"denom"`
// 				LazySchedules struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"lazy_schedules"`
// 			} `yaml:"properties"`
// 		} `yaml:"VestingSchedule"`
// 		BaseLazyGradedVestingAccount struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				BaseVestingAccount struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"BaseVestingAccount"`
// 				VestingSchedules struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"vesting_schedules"`
// 			} `yaml:"properties"`
// 		} `yaml:"BaseLazyGradedVestingAccount"`
// 		LazyGradedVestingAccount struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Type struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"type"`
// 				Value struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"value"`
// 			} `yaml:"properties"`
// 		} `yaml:"LazyGradedVestingAccount"`
// 		SwapReq struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				BaseReq struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"base_req"`
// 				OfferCoin struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"offer_coin"`
// 				AskDenom struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"ask_denom"`
// 			} `yaml:"properties"`
// 		} `yaml:"SwapReq"`
// 		MarketParams struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				DailyLunaDeltaLimit struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"daily_luna_delta_limit"`
// 				MinSwapSpread struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"min_swap_spread"`
// 				MaxSwapSpread struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"max_swap_spread"`
// 			} `yaml:"properties"`
// 		} `yaml:"MarketParams"`
// 		PrevoteReq struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				BaseReq struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"base_req"`
// 				ExchangeRate struct {
// 					Type        string `yaml:"type"`
// 					Example     string `yaml:"example"`
// 					Description string `yaml:"description"`
// 				} `yaml:"exchange_rate"`
// 				Salt struct {
// 					Type        string `yaml:"type"`
// 					Example     string `yaml:"example"`
// 					Description string `yaml:"description"`
// 				} `yaml:"salt"`
// 				Hash struct {
// 					Type        string `yaml:"type"`
// 					Example     string `yaml:"example"`
// 					Description string `yaml:"description"`
// 				} `yaml:"hash"`
// 				Validator struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"validator"`
// 			} `yaml:"properties"`
// 		} `yaml:"PrevoteReq"`
// 		VoteReq struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				BaseReq struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"base_req"`
// 				ExchangeRate struct {
// 					Type        string `yaml:"type"`
// 					Example     string `yaml:"example"`
// 					Description string `yaml:"description"`
// 				} `yaml:"exchange_rate"`
// 				Salt struct {
// 					Type        string `yaml:"type"`
// 					Example     string `yaml:"example"`
// 					Description string `yaml:"description"`
// 				} `yaml:"salt"`
// 				Validator struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"validator"`
// 			} `yaml:"properties"`
// 		} `yaml:"VoteReq"`
// 		DelegateReq struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				BaseReq struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"base_req"`
// 				Feeder struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"feeder"`
// 			} `yaml:"properties"`
// 		} `yaml:"DelegateReq"`
// 		ExchangeRateVote struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				ExchangeRate struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"exchange_rate"`
// 				Denom struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"denom"`
// 				Voter struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"voter"`
// 			} `yaml:"properties"`
// 		} `yaml:"ExchangeRateVote"`
// 		ExchangeRatePrevote struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Hash struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"hash"`
// 				Denom struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"denom"`
// 				Voter struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"voter"`
// 				SubmitBlock struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"submit_block"`
// 			} `yaml:"properties"`
// 		} `yaml:"ExchangeRatePrevote"`
// 		OracleParams struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				VotePeriod struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"vote_period"`
// 				VoteThreshold struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"vote_threshold"`
// 				DropThreshold struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"drop_threshold"`
// 				OracleRewardBand struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"oracle_reward_band"`
// 			} `yaml:"properties"`
// 		} `yaml:"OracleParams"`
// 		PolicyConstraints struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				RateMin struct {
// 					Type        string `yaml:"type"`
// 					Format      string `yaml:"format"`
// 					Example     string `yaml:"example"`
// 					Description string `yaml:"description"`
// 				} `yaml:"rate_min"`
// 				RateMax struct {
// 					Type        string `yaml:"type"`
// 					Format      string `yaml:"format"`
// 					Example     string `yaml:"example"`
// 					Description string `yaml:"description"`
// 				} `yaml:"rate_max"`
// 				Cap struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"cap"`
// 				ChangeMax struct {
// 					Type        string `yaml:"type"`
// 					Format      string `yaml:"format"`
// 					Example     string `yaml:"example"`
// 					Description string `yaml:"description"`
// 				} `yaml:"change_max"`
// 			} `yaml:"properties"`
// 		} `yaml:"PolicyConstraints"`
// 		TreasuryParams struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				TaxPolicy struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"tax_policy"`
// 				RewardPolicy struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"reward_policy"`
// 				SeigniorageBurdenTarget struct {
// 					Type        string `yaml:"type"`
// 					Format      string `yaml:"format"`
// 					Example     string `yaml:"example"`
// 					Description string `yaml:"description"`
// 				} `yaml:"seigniorage_burden_target"`
// 				MiningIncrement struct {
// 					Type    string `yaml:"type"`
// 					Format  string `yaml:"format"`
// 					Example string `yaml:"example"`
// 				} `yaml:"mining_increment"`
// 				WindowShort struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"window_short"`
// 				WindowLong struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"window_long"`
// 				WindowProbation struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"window_probation"`
// 				OracleShare struct {
// 					Type    string `yaml:"type"`
// 					Format  string `yaml:"format"`
// 					Example string `yaml:"example"`
// 				} `yaml:"oracle_share"`
// 				BudgetShare struct {
// 					Type    string `yaml:"type"`
// 					Format  string `yaml:"format"`
// 					Example string `yaml:"example"`
// 				} `yaml:"budget_share"`
// 			} `yaml:"properties"`
// 		} `yaml:"TreasuryParams"`
// 		MultiSignPubKey struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Threshold struct {
// 					Type    string `yaml:"type"`
// 					Example int    `yaml:"example"`
// 				} `yaml:"threshold"`
// 				Pubkeys struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Type    string `yaml:"type"`
// 						Example string `yaml:"example"`
// 					} `yaml:"items"`
// 				} `yaml:"pubkeys"`
// 			} `yaml:"properties"`
// 		} `yaml:"MultiSignPubKey"`
// 		MultiSignReq struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Tx struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"tx"`
// 				ChainID struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"chain_id"`
// 				Signatures struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"signatures"`
// 				SignatureOnly struct {
// 					Type string `yaml:"type"`
// 				} `yaml:"signature_only"`
// 				Pubkey struct {
// 					Ref string `yaml:"$ref"`
// 				} `yaml:"pubkey"`
// 			} `yaml:"properties"`
// 		} `yaml:"MultiSignReq"`
// 		StdSignature struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Signature struct {
// 					Type    string `yaml:"type"`
// 					Example string `yaml:"example"`
// 				} `yaml:"signature"`
// 				PubKey struct {
// 					Type       string `yaml:"type"`
// 					Properties struct {
// 						Type struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"type"`
// 						Value struct {
// 							Type    string `yaml:"type"`
// 							Example string `yaml:"example"`
// 						} `yaml:"value"`
// 					} `yaml:"properties"`
// 				} `yaml:"pub_key"`
// 			} `yaml:"properties"`
// 		} `yaml:"StdSignature"`
// 		EstimateFeeResp struct {
// 			Type       string `yaml:"type"`
// 			Properties struct {
// 				Fees struct {
// 					Type  string `yaml:"type"`
// 					Items struct {
// 						Ref string `yaml:"$ref"`
// 					} `yaml:"items"`
// 				} `yaml:"fees"`
// 				Gas struct {
// 					Type    string `yaml:"type"`
// 					Format  string `yaml:"format"`
// 					Example string `yaml:"example"`
// 				} `yaml:"gas"`
// 			} `yaml:"properties"`
// 		} `yaml:"EstimateFeeResp"`
// 	} `yaml:"definitions"`
// }

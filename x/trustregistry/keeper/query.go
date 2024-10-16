package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"errors"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sort"
)

var _ types.QueryServer = queryServer{}

func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

func (qs queryServer) GetTrustRegistry(ctx context.Context, req *types.QueryGetTrustRegistryRequest) (*types.QueryGetTrustRegistryResponse, error) {
	if !isValidDID(req.Did) {
		return nil, status.Error(codes.InvalidArgument, "invalid DID syntax")
	}

	tr, err := qs.k.TrustRegistry.Get(ctx, req.Did)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "trust registry not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	var versions []types.GovernanceFrameworkVersion
	var documents []types.GovernanceFrameworkDocument

	// Fetch versions
	err = qs.k.GFVersion.Walk(ctx, nil, func(key string, gfv types.GovernanceFrameworkVersion) (bool, error) {
		if gfv.TrDid == req.Did && (!req.ActiveGfOnly || gfv.Version == tr.ActiveVersion) {
			versions = append(versions, gfv)
		}
		return false, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Fetch documents
	err = qs.k.GFDocument.Walk(ctx, nil, func(key string, gfd types.GovernanceFrameworkDocument) (bool, error) {
		for _, v := range versions {
			if gfd.GfvId == v.Id {
				if req.PreferredLanguage == "" || gfd.Language == req.PreferredLanguage {
					documents = append(documents, gfd)
					break
				} else if len(documents) == 0 || documents[len(documents)-1].GfvId != v.Id {
					documents = append(documents, gfd)
				}
			}
		}
		return false, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetTrustRegistryResponse{
		TrustRegistry: &tr,
		Versions:      versions,
		Documents:     documents,
	}, nil
}

func (qs queryServer) ListTrustRegistries(ctx context.Context, req *types.QueryListTrustRegistriesRequest) (*types.QueryListTrustRegistriesResponse, error) {
	if req.ResponseMaxSize < 1 || req.ResponseMaxSize > 1024 {
		return nil, status.Error(codes.InvalidArgument, "response_max_size must be between 1 and 1024")
	}

	var trustRegistries []types.TrustRegistry
	var err error

	err = qs.k.TrustRegistry.Walk(ctx, nil, func(key string, tr types.TrustRegistry) (bool, error) {
		if req.Modified != nil && !tr.Modified.After(*req.Modified) {
			return false, nil
		}
		trustRegistries = append(trustRegistries, tr)
		return len(trustRegistries) >= int(req.ResponseMaxSize), nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Sort by modified time ascending
	sort.Slice(trustRegistries, func(i, j int) bool {
		return trustRegistries[i].Modified.Before(trustRegistries[j].Modified)
	})

	return &types.QueryListTrustRegistriesResponse{
		TrustRegistries: trustRegistries,
	}, nil
}

// Params defines the handler for the Query/Params RPC method.
func (qs queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := qs.k.Params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &types.QueryParamsResponse{Params: types.Params{}}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

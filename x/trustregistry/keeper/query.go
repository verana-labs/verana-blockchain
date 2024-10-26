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
	if req.TrId == 0 {
		return nil, status.Error(codes.InvalidArgument, "trust registry ID is required")
	}

	// Direct lookup by ID
	tr, err := qs.k.TrustRegistry.Get(ctx, req.TrId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "trust registry not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return qs.getTrustRegistryData(ctx, tr, req.ActiveGfOnly, req.PreferredLanguage)
}

func (qs queryServer) GetTrustRegistryWithDID(ctx context.Context, req *types.QueryGetTrustRegistryWithDIDRequest) (*types.QueryGetTrustRegistryResponse, error) {
	if !isValidDID(req.Did) {
		return nil, status.Error(codes.InvalidArgument, "invalid DID syntax")
	}

	// Get ID from DID index
	id, err := qs.k.TrustRegistryDIDIndex.Get(ctx, req.Did)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "trust registry not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Get trust registry using ID
	tr, err := qs.k.TrustRegistry.Get(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return qs.getTrustRegistryData(ctx, tr, req.ActiveGfOnly, req.PreferredLanguage)
}

func (qs queryServer) getTrustRegistryData(ctx context.Context, tr types.TrustRegistry, activeOnly bool, preferredLang string) (*types.QueryGetTrustRegistryResponse, error) {
	var versions []types.GovernanceFrameworkVersion
	var documents []types.GovernanceFrameworkDocument

	// Fetch versions
	err := qs.k.GFVersion.Walk(ctx, nil, func(id uint64, gfv types.GovernanceFrameworkVersion) (bool, error) {
		if gfv.TrId == tr.Id {
			if !activeOnly || gfv.Version == tr.ActiveVersion {
				versions = append(versions, gfv)
			}
		}
		return false, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Fetch documents
	for _, v := range versions {
		var versionDocs []types.GovernanceFrameworkDocument
		err = qs.k.GFDocument.Walk(ctx, nil, func(id uint64, gfd types.GovernanceFrameworkDocument) (bool, error) {
			if gfd.GfvId == v.Id {
				versionDocs = append(versionDocs, gfd)
			}
			return false, nil
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		// If preferred language is set, try to find a matching document
		if preferredLang != "" {
			var preferredDoc *types.GovernanceFrameworkDocument
			var fallbackDoc *types.GovernanceFrameworkDocument

			for i, doc := range versionDocs {
				if doc.Language == preferredLang {
					preferredDoc = &versionDocs[i]
					break
				} else if fallbackDoc == nil {
					fallbackDoc = &versionDocs[i]
				}
			}

			// Add preferred language doc if found, otherwise add fallback
			if preferredDoc != nil {
				documents = append(documents, *preferredDoc)
			} else if fallbackDoc != nil {
				documents = append(documents, *fallbackDoc)
			}
		} else {
			// If no preferred language, add all documents
			documents = append(documents, versionDocs...)
		}
	}

	return &types.QueryGetTrustRegistryResponse{
		TrustRegistry: &tr,
		Versions:      versions,
		Documents:     documents,
	}, nil
}

func (qs queryServer) ListTrustRegistries(ctx context.Context, req *types.QueryListTrustRegistriesRequest) (*types.QueryListTrustRegistriesResponse, error) {
	// Validate response_max_size
	if req.ResponseMaxSize < 1 || req.ResponseMaxSize > 1024 {
		return nil, status.Error(codes.InvalidArgument, "response_max_size must be between 1 and 1,024")
	}

	var trustRegistries []types.TrustRegistry

	// Collect all matching trust registries
	err := qs.k.TrustRegistry.Walk(ctx, nil, func(key uint64, tr types.TrustRegistry) (bool, error) {
		// Apply filters
		if req.Controller != "" && tr.Controller != req.Controller {
			return false, nil
		}
		if req.ModifiedAfter != nil && !tr.Modified.After(*req.ModifiedAfter) {
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

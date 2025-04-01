package retag

import (
	"github.com/spf13/cobra"
)

const (
	name             = "retag"
	shortDescription = "Retags image tags specified in the configuration file"
	longDescription  = `Retag will iterate the specified configuration file and Retag only the image tags listed.
If the image tag is already in the registry it will update that tag if the digests have changed.`

	TimebombInMonths       = 6
	DescriptorFetchRetries = 6

	RetagAZ    = "az"
	RetagCrane = "oras"
	RetagOras  = "oras"
)

type retagCmd struct {
	source              string
	tags                []string
	registry            string
	destination         string
	dateAdded           string
	enableTimeBomb      bool
	dryrun              bool
	generateSupplyChain bool
	outputDir           string
	tool                string
}

func NewRetagCmd() *cobra.Command {
	//rc := &retagCmd{}
	retagCmd := &cobra.Command{
		Use:   name,
		Short: shortDescription,
		Long:  longDescription,
	}

	retagCmd.AddCommand(newGenerateCommand())
	return retagCmd
}

// func (rc *retagCmd) digestsEqual(source, dest string) (bool, string, string, error) {
// 	sourceDigest, err := rc.oras.GetDigest(context.TODO(), source)
// 	if err != nil {
// 		return false, "", "", err
// 	}
// 	destDigest, err := rc.oras.GetDigest(context.TODO(), dest)
// 	// return the error if it doesn't contain "not found" (when the ref doesn't exist in dest)
// 	if err != nil && !strings.Contains(err.Error(), "not found") {
// 		return false, "", "", err
// 	}
// 	return sourceDigest == destDigest, sourceDigest, destDigest, nil
// }

// func (rc *retagCmd) AttachEOLAndSign(ref string, tag string) error {
// 	digest, err := rc.oras.AttachMicrosoftEOLAnnotation(ref, nil)
// 	if err != nil {
// 		return fmt.Errorf("unable to attach EOL to %s (previously associated with %s): %s", ref, tag, err)
// 	}
// 	eolAnnotationsDir := filepath.Join(rc.outputDir, "retag_results", "eol_annotations")
// 	fullEOLRef := fmt.Sprintf("%s@%s", rc.fullRepoURL(), digest)

// 	err = rc.writeSigningDescriptor(eolAnnotationsDir, fullEOLRef)
// 	if err != nil {
// 		return fmt.Errorf("unable to attach signing descriptor for %s eol annotations: %s", digest, err)
// 	}

// 	return nil
// }

// func (rc *retagCmd) timebombExpired() (bool, error) {
// 	if !rc.enableTimeBomb {
// 		return false, nil
// 	}

// 	dateAdded, err := time.Parse("2006-01-02", rc.dateAdded)
// 	if err != nil {
// 		return false, errors.Wrapf(err, "failed to parse date: %s", rc.dateAdded)
// 	}

// 	if !dateAdded.AddDate(0, TimebombInMonths, 0).After(time.Now()) {
// 		return true, nil
// 	}

// 	return false, nil
// }

// func (rc *retagCmd) generateSupplyChainArtifacts(registry, repo, tag string) error {
// 	retagResultsDir := filepath.Join(rc.outputDir, "retag_results")
// 	dir := filepath.Join(retagResultsDir, "sboms")
// 	if _, err := os.Stat(dir); os.IsNotExist(err) {
// 		if err := os.MkdirAll(dir, 0o755); err != nil {
// 			return err
// 		}
// 	}
// 	ref := fmt.Sprintf("%s/%s:%s", registry, repo, tag)
// 	// this will get the descriptor for the manifest list
// 	descriptor, err := rc.oras.GetDescriptor(context.Background(), ref)
// 	if err != nil {
// 		return err
// 	}
// 	signingFileName := fmt.Sprintf("%s:%s", repo, tag)
// 	path, err := signing.WriteDescriptor(
// 		retagResultsDir,
// 		signingFileName,
// 		descriptor,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	log.Infof("Wrote a descriptor file for %s at: %s", ref, path)
// 	// get the manifest to determine if we need to handle this as a manifest list or a manifest
// 	manifest, err := rc.docker.InspectManifest(ref)
// 	if err != nil {
// 		return err
// 	}
// 	// if the mediaType is manifest, we need to generate an sbom and then return
// 	if manifest.MediaType == "application/vnd.docker.distribution.manifest.v2+json" {
// 		if err := rc.generateSBOM(dir, registry, repo, tag, ""); err != nil {
// 			log.Warnf("unable to generate sbom for %s:%s: %s", repo, tag, err)
// 		}
// 		return nil
// 	}
// 	// this is a manifest list, so use the digests to make references
// 	for _, manifest := range manifest.Manifests {
// 		shaRef := fmt.Sprintf("%s/%s@%s", registry, repo, manifest.Digest)
// 		imgDesc, err := rc.oras.GetDescriptor(context.Background(), shaRef)
// 		if err != nil {
// 			return err
// 		}
// 		signingFileName = fmt.Sprintf("%s@%s", repo, manifest.Digest)
// 		path, err := signing.WriteDescriptor(
// 			retagResultsDir,
// 			signingFileName,
// 			imgDesc,
// 		)
// 		if err != nil {
// 			return err
// 		}
// 		log.Infof("Wrote a descriptor file for %s at: %s", shaRef, path)
// 		// generate an SBOM for the image if it is linux
// 		if manifest.Platform.OS == "linux" {
// 			if err := rc.generateSBOM(dir, registry, repo, tag, manifest.Digest); err != nil {
// 				return fmt.Errorf("unable to generate SBOM for %s: %s", shaRef, err)
// 			}
// 		} else {
// 			platformOS := manifest.Platform.OS
// 			if platformOS == "" {
// 				platformOS = "unknown"
// 			}
// 			log.Infof("skipping SBOM generation for %s: OS=%s", shaRef, platformOS)
// 		}
// 	}
// 	return nil
// }

// func (rc *retagCmd) writeSigningDescriptor(outputDir string, digestRef string) error {
// 	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
// 		if err := os.MkdirAll(outputDir, 0o755); err != nil {
// 			return err
// 		}
// 	}
// 	splitDigestRef := strings.Split(digestRef, "/")
// 	signingRef := splitDigestRef[len(splitDigestRef)-1]
// 	desc, err := rc.oras.GetDescriptor(context.Background(), digestRef)
// 	if err != nil {
// 		return fmt.Errorf("unable to get descriptor for %s: %s", digestRef, err)
// 	}
// 	path, err := signing.WriteDescriptor(outputDir, signingRef, desc)
// 	if err != nil {
// 		return fmt.Errorf("unable to write signing descriptor for %s: %s", signingRef, err)
// 	}
// 	log.Infof("wrote signing descriptor for %s to %s", signingRef, path)
// 	return nil
// }

// func (rc *retagCmd) generateSBOM(dir, registry, repo, tag, digest string) error {
// 	imageRef := fmt.Sprintf("%s/%s:%s", registry, repo, tag)
// 	if digest != "" {
// 		imageRef = fmt.Sprintf("%s/%s@%s", registry, repo, digest)
// 	}
// 	lastSlashPosition := strings.LastIndex(repo, "/")
// 	imageName := repo[lastSlashPosition:]
// 	sbomFileName := fmt.Sprintf("%s:%s", imageName, tag)
// 	if digest != "" {
// 		sbomFileName = fmt.Sprintf("%s:%s@%s", imageName, tag, digest)
// 	}
// 	fileName := filepath.Join(dir, sbomFileName)
// 	log.Infof("Generating SBOM at %s for %s", fileName, imageRef)
// 	if err := rc.syft.Packages(imageRef, "spdx-json", fileName); err != nil {
// 		log.Warnf("unable to write sbom for %s", imageRef)
// 	}
// 	return nil
// }

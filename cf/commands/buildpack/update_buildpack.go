package buildpack

import (
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/command_metadata"
	"github.com/cloudfoundry/cli/cf/flag_helpers"
	. "github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type UpdateBuildpack struct {
	ui                terminal.UI
	buildpackRepo     api.BuildpackRepository
	buildpackBitsRepo api.BuildpackBitsRepository
	buildpackReq      requirements.BuildpackRequirement
}

func NewUpdateBuildpack(ui terminal.UI, repo api.BuildpackRepository, bitsRepo api.BuildpackBitsRepository) (cmd *UpdateBuildpack) {
	cmd = new(UpdateBuildpack)
	cmd.ui = ui
	cmd.buildpackRepo = repo
	cmd.buildpackBitsRepo = bitsRepo
	return
}

func (cmd *UpdateBuildpack) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "update-buildpack",
		Description: T("Update a buildpack"),
		Usage: T("CF_NAME update-buildpack BUILDPACK [-p PATH] [-i POSITION] [--enable|--disable] [--lock|--unlock]") +
			T("\n\nTIP:\n") + T("   Path should be a zip file, a url to a zip file, or a local directory. Position is a positive integer, sets priority, and is sorted from lowest to highest."),
		Flags: []cli.Flag{
			flag_helpers.NewIntFlag("i", T("Buildpack position among other buildpacks")),
			flag_helpers.NewStringFlag("p", T("Path to directory or zip file")),
			cli.BoolFlag{Name: "enable", Usage: T("Enable the buildpack")},
			cli.BoolFlag{Name: "disable", Usage: T("Disable the buildpack")},
			cli.BoolFlag{Name: "lock", Usage: T("Lock the buildpack")},
			cli.BoolFlag{Name: "unlock", Usage: T("Unlock the buildpack")},
		},
	}
}

func (cmd *UpdateBuildpack) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
		cmd.ui.FailWithUsage(c)
	}

	loginReq := requirementsFactory.NewLoginRequirement()
	cmd.buildpackReq = requirementsFactory.NewBuildpackRequirement(c.Args()[0])

	reqs = []requirements.Requirement{
		loginReq,
		cmd.buildpackReq,
	}

	return
}

func (cmd *UpdateBuildpack) Run(c *cli.Context) {
	buildpack := cmd.buildpackReq.GetBuildpack()

	cmd.ui.Say(T("Updating buildpack {{.BuildpackName}}...", map[string]interface{}{"BuildpackName": terminal.EntityNameColor(buildpack.Name)}))

	updateBuildpack := false

	if c.IsSet("i") {
		position := c.Int("i")

		buildpack.Position = &position
		updateBuildpack = true
	}

	enabled := c.Bool("enable")
	disabled := c.Bool("disable")
	if enabled && disabled {
		cmd.ui.Failed(T("Cannot specify both {{.Enabled}} and {{.Disabled}}.", map[string]interface{}{
			"Enabled":  "enabled",
			"Disabled": "disabled",
		}))
	}

	if enabled {
		buildpack.Enabled = &enabled
		updateBuildpack = true
	}
	if disabled {
		disabled = false
		buildpack.Enabled = &disabled
		updateBuildpack = true
	}

	lock := c.Bool("lock")
	unlock := c.Bool("unlock")
	if lock && unlock {
		cmd.ui.Failed(T("Cannot specify both lock and unlock options."))
		return
	}

	dir := c.String("p")
	if dir != "" && (lock || unlock) {
		cmd.ui.Failed(T("Cannot specify buildpack bits and lock/unlock."))
	}

	if lock {
		buildpack.Locked = &lock
		updateBuildpack = true
	}
	if unlock {
		unlock = false
		buildpack.Locked = &unlock
		updateBuildpack = true
	}

	if updateBuildpack {
		newBuildpack, apiErr := cmd.buildpackRepo.Update(buildpack)
		if apiErr != nil {
			cmd.ui.Failed(T("Error updating buildpack {{.Name}}\n{{.Error}}", map[string]interface{}{
				"Name":  terminal.EntityNameColor(buildpack.Name),
				"Error": apiErr.Error(),
			}))
		}
		buildpack = newBuildpack
	}

	if dir != "" {
		apiErr := cmd.buildpackBitsRepo.UploadBuildpack(buildpack, dir)
		if apiErr != nil {
			cmd.ui.Failed(T("Error uploading buildpack {{.Name}}\n{{.Error}}", map[string]interface{}{
				"Name":  terminal.EntityNameColor(buildpack.Name),
				"Error": apiErr.Error(),
			}))
		}
	}
	cmd.ui.Ok()
}

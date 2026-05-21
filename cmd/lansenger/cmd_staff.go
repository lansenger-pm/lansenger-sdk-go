package main

import (
	"context"

	"github.com/spf13/cobra"
)

var staffCmd = &cobra.Command{
	Use:   "staff",
	Short: "Staff and contact commands",
}

var staffBasicInfoCmd = &cobra.Command{
	Use:   "basic-info STAFF_ID",
	Short: "Fetch staff basic info",
	Args:  cobra.ExactArgs(1),
	Run:   runStaffBasicInfo,
}

var staffDetailCmd = &cobra.Command{
	Use:   "detail STAFF_ID",
	Short: "Fetch staff detail info",
	Args:  cobra.ExactArgs(1),
	Run:   runStaffDetail,
}

var staffAncestorsCmd = &cobra.Command{
	Use:   "ancestors STAFF_ID",
	Short: "Fetch department ancestors for a staff",
	Args:  cobra.ExactArgs(1),
	Run:   runStaffAncestors,
}

var staffIdMappingCmd = &cobra.Command{
	Use:   "id-mapping ORG_ID ID_TYPE ID_VALUE",
	Short: "Fetch staff ID mapping",
	Args:  cobra.ExactArgs(3),
	Run:   runStaffIdMapping,
}

var staffOrgExtraFieldsCmd = &cobra.Command{
	Use:   "org-extra-fields ORG_ID",
	Short: "Fetch org extra field IDs",
	Args:  cobra.ExactArgs(1),
	Run:   runStaffOrgExtraFields,
}

var staffSearchCmd = &cobra.Command{
	Use:   "search KEYWORD",
	Short: "Search staff",
	Args:  cobra.ExactArgs(1),
	Run:   runStaffSearch,
}

var staffOrgInfoCmd = &cobra.Command{
	Use:   "org-info ORG_ID",
	Short: "Fetch org info",
	Args:  cobra.ExactArgs(1),
	Run:   runStaffOrgInfo,
}

var (
	staffBasicInfoUserToken  string
	staffDetailUserToken     string
	staffAncestorsUserToken  string
	staffIdMappingUserToken  string
	staffOrgExtraFieldsToken string
	staffOrgExtraFieldsPage  int
	staffOrgExtraFieldsSize  int
	staffSearchUserToken     string
	staffSearchUserID        string
	staffSearchRecursive     bool
	staffSearchSectorIDs     []string
	staffSearchPage          int
	staffSearchSize          int
	staffOrgInfoUserToken    string
)

func init() {
	staffBasicInfoCmd.Flags().StringVar(&staffBasicInfoUserToken, "user-token", "", "User token")
	staffDetailCmd.Flags().StringVar(&staffDetailUserToken, "user-token", "", "User token")
	staffAncestorsCmd.Flags().StringVar(&staffAncestorsUserToken, "user-token", "", "User token")
	staffIdMappingCmd.Flags().StringVar(&staffIdMappingUserToken, "user-token", "", "User token")
	staffOrgExtraFieldsCmd.Flags().StringVar(&staffOrgExtraFieldsToken, "user-token", "", "User token")
	staffOrgExtraFieldsCmd.Flags().IntVarP(&staffOrgExtraFieldsPage, "page", "p", 1, "Page number")
	staffOrgExtraFieldsCmd.Flags().IntVarP(&staffOrgExtraFieldsSize, "size", "s", 20, "Page size")
	staffSearchCmd.Flags().StringVar(&staffSearchUserToken, "user-token", "", "User token")
	staffSearchCmd.Flags().StringVar(&staffSearchUserID, "user-id", "", "User ID")
	staffSearchCmd.Flags().BoolVar(&staffSearchRecursive, "recursive", true, "Recursive search")
	staffSearchCmd.Flags().StringArrayVar(&staffSearchSectorIDs, "sector", nil, "Sector IDs")
	staffSearchCmd.Flags().IntVarP(&staffSearchPage, "page", "p", 1, "Page number")
	staffSearchCmd.Flags().IntVarP(&staffSearchSize, "size", "s", 20, "Page size")
	staffOrgInfoCmd.Flags().StringVar(&staffOrgInfoUserToken, "user-token", "", "User token")

	staffCmd.AddCommand(staffBasicInfoCmd)
	staffCmd.AddCommand(staffDetailCmd)
	staffCmd.AddCommand(staffAncestorsCmd)
	staffCmd.AddCommand(staffIdMappingCmd)
	staffCmd.AddCommand(staffOrgExtraFieldsCmd)
	staffCmd.AddCommand(staffSearchCmd)
	staffCmd.AddCommand(staffOrgInfoCmd)
	rootCmd.AddCommand(staffCmd)
}

func runStaffBasicInfo(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchStaffBasicInfo(ctx, args[0], staffBasicInfoUserToken)
	checkError(err)
	outputResultFields(result, []string{"org_id", "org_name", "name", "gender", "signature", "avatar_url", "status", "departments"})
}

func runStaffDetail(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchStaffDetail(ctx, args[0], staffDetailUserToken)
	checkError(err)
	outputResultFields(result, []string{"org_id", "org_name", "name", "gender", "email", "mobile_phone", "avatar_url", "career", "tags"})
}

func runStaffAncestors(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchDepartmentAncestors(ctx, args[0], staffAncestorsUserToken)
	checkError(err)
	outputResult(result)
}

func runStaffIdMapping(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchStaffIdMapping(ctx, args[0], args[1], args[2], staffIdMappingUserToken)
	checkError(err)
	outputResultFields(result, []string{"staff_id"})
}

func runStaffOrgExtraFields(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchOrgExtraFieldIDs(ctx, args[0], staffOrgExtraFieldsToken, staffOrgExtraFieldsPage, staffOrgExtraFieldsSize)
	checkError(err)
	outputResultFields(result, []string{"has_more", "total", "extra_field_ids"})
}

func runStaffSearch(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	var sectorIDs []string
	if len(staffSearchSectorIDs) > 0 {
		sectorIDs = staffSearchSectorIDs
	}

	result, err := client.SearchStaff(ctx, args[0], staffSearchUserToken, staffSearchUserID, staffSearchRecursive, sectorIDs, staffSearchPage, staffSearchSize)
	checkError(err)
	outputResultFields(result, []string{"has_more", "total", "staff_info"})
}

func runStaffOrgInfo(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchOrgInfo(ctx, args[0], staffOrgInfoUserToken)
	checkError(err)
	outputResultFields(result, []string{"org_id", "org_name", "icon_url"})
}
package main

import (
	"context"
	"strconv"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Group commands",
}

var groupCreateCmd = &cobra.Command{
	Use:   "create NAME ORG_ID",
	Short: "Create a group",
	Args:  cobra.ExactArgs(2),
	Run:   runGroupCreate,
}

var groupInfoCmd = &cobra.Command{
	Use:   "info GROUP_ID",
	Short: "Fetch group info",
	Args:  cobra.ExactArgs(1),
	Run:   runGroupInfo,
}

var groupMembersCmd = &cobra.Command{
	Use:   "members GROUP_ID",
	Short: "Fetch group members",
	Args:  cobra.ExactArgs(1),
	Run:   runGroupMembers,
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "Fetch group list",
	Args:  cobra.NoArgs,
	Run:   runGroupList,
}

var groupCheckCmd = &cobra.Command{
	Use:   "check GROUP_ID",
	Short: "Check if staff is in group",
	Args:  cobra.ExactArgs(1),
	Run:   runGroupCheck,
}

var groupUpdateCmd = &cobra.Command{
	Use:   "update GROUP_ID",
	Short: "Update group info",
	Args:  cobra.ExactArgs(1),
	Run:   runGroupUpdate,
}

var groupUpdateMembersCmd = &cobra.Command{
	Use:   "update-members GROUP_ID",
	Short: "Update group members",
	Args:  cobra.ExactArgs(1),
	Run:   runGroupUpdateMembers,
}

var groupDismissCmd = &cobra.Command{
	Use:   "dismiss GROUP_ID",
	Short: "Dismiss (delete) a group",
	Args:  cobra.ExactArgs(1),
	Run:   runGroupDismiss,
}

var (
	groupCreateOwner      string
	groupCreateDesc       string
	groupCreateAvatar     string
	groupCreateStaff      []string
	groupCreateDept       []string
	groupCreateUserToken  string

	groupInfoUserToken string

	groupMembersUserToken string
	groupMembersPage      int
	groupMembersSize      int

	groupListUserToken string
	groupListPage      int
	groupListSize      int

	groupCheckUserToken string
	groupCheckStaffID   string

	groupUpdateName      string
	groupUpdateDesc      string
	groupUpdateOwner     string
	groupUpdateUserToken string

	groupUpdateMembersAdd    []string
	groupUpdateMembersRemove []string
	groupUpdateMembersAddDept []string
	groupUpdateMembersToken  string

	groupDismissUserToken string
)

func init() {
	groupCreateCmd.Flags().StringVar(&groupCreateOwner, "owner", "", "Owner ID")
	groupCreateCmd.Flags().StringVarP(&groupCreateDesc, "desc", "d", "", "Description")
	groupCreateCmd.Flags().StringVar(&groupCreateAvatar, "avatar", "", "Avatar ID")
	groupCreateCmd.Flags().StringArrayVar(&groupCreateStaff, "staff", nil, "Staff IDs (repeatable)")
	groupCreateCmd.Flags().StringArrayVar(&groupCreateDept, "dept", nil, "Department IDs (repeatable)")
	groupCreateCmd.Flags().StringVar(&groupCreateUserToken, "user-token", "", "User token")

	groupInfoCmd.Flags().StringVar(&groupInfoUserToken, "user-token", "", "User token")

	groupMembersCmd.Flags().StringVar(&groupMembersUserToken, "user-token", "", "User token")
	groupMembersCmd.Flags().IntVarP(&groupMembersPage, "page", "p", 0, "Page offset")
	groupMembersCmd.Flags().IntVarP(&groupMembersSize, "size", "s", 100, "Page size")

	groupListCmd.Flags().StringVar(&groupListUserToken, "user-token", "", "User token")
	groupListCmd.Flags().IntVarP(&groupListPage, "page", "p", 0, "Page offset")
	groupListCmd.Flags().IntVarP(&groupListSize, "size", "s", 100, "Page size")

	groupCheckCmd.Flags().StringVar(&groupCheckUserToken, "user-token", "", "User token")
	groupCheckCmd.Flags().StringVar(&groupCheckStaffID, "staff-id", "", "Staff ID to check")

	groupUpdateCmd.Flags().StringVar(&groupUpdateName, "name", "", "New group name")
	groupUpdateCmd.Flags().StringVar(&groupUpdateDesc, "desc", "", "New description")
	groupUpdateCmd.Flags().StringVar(&groupUpdateOwner, "owner", "", "New owner ID")
	groupUpdateCmd.Flags().StringVar(&groupUpdateUserToken, "user-token", "", "User token")

	groupUpdateMembersCmd.Flags().StringArrayVar(&groupUpdateMembersAdd, "add", nil, "Staff IDs to add (repeatable)")
	groupUpdateMembersCmd.Flags().StringArrayVar(&groupUpdateMembersRemove, "remove", nil, "Staff IDs to remove (repeatable)")
	groupUpdateMembersCmd.Flags().StringArrayVar(&groupUpdateMembersAddDept, "add-dept", nil, "Department IDs to add (repeatable)")
	groupUpdateMembersCmd.Flags().StringVar(&groupUpdateMembersToken, "user-token", "", "User token")

	groupDismissCmd.Flags().StringVar(&groupDismissUserToken, "user-token", "", "User token")

	groupCmd.AddCommand(groupCreateCmd)
	groupCmd.AddCommand(groupInfoCmd)
	groupCmd.AddCommand(groupMembersCmd)
	groupCmd.AddCommand(groupListCmd)
	groupCmd.AddCommand(groupCheckCmd)
	groupCmd.AddCommand(groupUpdateCmd)
	groupCmd.AddCommand(groupUpdateMembersCmd)
	groupCmd.AddCommand(groupDismissCmd)
	rootCmd.AddCommand(groupCmd)
}

func runGroupCreate(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	orgIDInt, err := strconv.Atoi(args[1])
	checkError(err)

	var staffIDList []string
	if len(groupCreateStaff) > 0 {
		staffIDList = groupCreateStaff
	}

	var deptIDList []string
	if len(groupCreateDept) > 0 {
		deptIDList = groupCreateDept
	}

	info := &lansenger.GroupCreateInfo{
		Name:             args[0],
		OrgID:            orgIDInt,
		OwnerID:          groupCreateOwner,
		Description:      groupCreateDesc,
		AvatarID:         groupCreateAvatar,
		StaffIDList:      staffIDList,
		DepartmentIDList: deptIDList,
	}

	result, err := client.CreateGroup(ctx, info, groupCreateUserToken)
	checkError(err)
	outputResultFields(result, []string{"group_id", "total_members"})
}

func runGroupInfo(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchGroupInfo(ctx, args[0], groupInfoUserToken)
	checkError(err)
	outputResult(result)
}

func runGroupMembers(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchGroupMembers(ctx, args[0], groupMembersUserToken, groupMembersPage, groupMembersSize)
	checkError(err)
	outputResultFields(result, []string{"total_members"})
}

func runGroupList(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchGroupList(ctx, groupListUserToken, groupListPage, groupListSize)
	checkError(err)
	outputResultFields(result, []string{"total_group_ids", "group_ids"})
}

func runGroupCheck(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.CheckIsInGroup(ctx, args[0], groupCheckUserToken, groupCheckStaffID)
	checkError(err)
	outputResultFields(result, []string{"is_in_group"})
}

func runGroupUpdate(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	params := map[string]interface{}{}
	if groupUpdateName != "" {
		params["name"] = groupUpdateName
	}
	if groupUpdateDesc != "" {
		params["description"] = groupUpdateDesc
	}
	if groupUpdateOwner != "" {
		params["ownerId"] = groupUpdateOwner
	}

	result, err := client.UpdateGroupInfo(ctx, args[0], params, groupUpdateUserToken)
	checkError(err)
	outputResult(result)
}

func runGroupDismiss(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.DissolveGroup(ctx, args[0], groupDismissUserToken)
	checkError(err)
	outputResult(result)
}

func runGroupUpdateMembers(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	var addUserList []string
	if len(groupUpdateMembersAdd) > 0 {
		addUserList = groupUpdateMembersAdd
	}

	var delUserList []string
	if len(groupUpdateMembersRemove) > 0 {
		delUserList = groupUpdateMembersRemove
	}

	var addDeptIDList []string
	if len(groupUpdateMembersAddDept) > 0 {
		addDeptIDList = groupUpdateMembersAddDept
	}

	result, err := client.UpdateGroupMembers(ctx, args[0], addUserList, delUserList, addDeptIDList, groupUpdateMembersToken)
	checkError(err)
	outputResultFields(result, []string{"total_members", "added_staff_count", "deleted_staff_count"})
}
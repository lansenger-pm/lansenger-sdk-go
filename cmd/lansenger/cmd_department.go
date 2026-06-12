package main

import (
	"context"

	"github.com/spf13/cobra"
)

var departmentCmd = &cobra.Command{
	Use:   "department",
	Short: "Query department information",
}

var departmentDetailCmd = &cobra.Command{
	Use:   "detail DEPARTMENT_ID",
	Short: "Fetch department detail",
	Args:  cobra.ExactArgs(1),
	Run:   runDepartmentDetail,
}

var departmentChildrenCmd = &cobra.Command{
	Use:   "children DEPARTMENT_ID",
	Short: "Fetch department children",
	Args:  cobra.ExactArgs(1),
	Run:   runDepartmentChildren,
}

var departmentStaffsCmd = &cobra.Command{
	Use:   "staffs DEPARTMENT_ID",
	Short: "Fetch department staffs",
	Args:  cobra.ExactArgs(1),
	Run:   runDepartmentStaffs,
}

var (
	departmentDetailUserToken string
	departmentDetailTagID     string
	departmentChildrenToken   string
	departmentStaffsToken     string
	departmentStaffsPage      int
	departmentStaffsSize      int
)

func init() {
	departmentDetailCmd.Flags().StringVar(&departmentDetailUserToken, "user-token", "", "User token")
	departmentDetailCmd.Flags().StringVar(&departmentDetailTagID, "tag-id", "", "Tag ID")
	departmentChildrenCmd.Flags().StringVar(&departmentChildrenToken, "user-token", "", "User token")
	departmentStaffsCmd.Flags().StringVar(&departmentStaffsToken, "user-token", "", "User token")
	departmentStaffsCmd.Flags().IntVarP(&departmentStaffsPage, "page", "p", 1, "Page number")
	departmentStaffsCmd.Flags().IntVarP(&departmentStaffsSize, "size", "s", 100, "Page size")

	departmentCmd.AddCommand(departmentDetailCmd)
	departmentCmd.AddCommand(departmentChildrenCmd)
	departmentCmd.AddCommand(departmentStaffsCmd)
	rootCmd.AddCommand(departmentCmd)
}

func runDepartmentDetail(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchDepartmentDetail(ctx, args[0], departmentDetailUserToken, departmentDetailTagID)
	checkError(err)
	outputResult(result)
}

func runDepartmentChildren(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchDepartmentChildren(ctx, args[0], departmentChildrenToken)
	checkError(err)
	outputResult(result)
}

func runDepartmentStaffs(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchDepartmentStaffs(ctx, args[0], departmentStaffsToken, departmentStaffsPage, departmentStaffsSize)
	checkError(err)
	outputResultFields(result, []string{"has_more", "total", "staffs"})
}
package services

import (
	"github.com/radhikadarode/finance-backend/internal/database"
	"github.com/radhikadarode/finance-backend/internal/models"
)

type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

func (s *DashboardService) GetSummary() (*models.DashboardSummary, error) {
	summary := &models.DashboardSummary{}

	// Total income
	database.DB.Model(&models.FinancialRecord{}).
		Where("type = ? AND is_deleted = ?", models.TypeIncome, false).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&summary.TotalIncome)

	// Total expenses
	database.DB.Model(&models.FinancialRecord{}).
		Where("type = ? AND is_deleted = ?", models.TypeExpense, false).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&summary.TotalExpenses)

	summary.NetBalance = summary.TotalIncome - summary.TotalExpenses

	// Total record count
	database.DB.Model(&models.FinancialRecord{}).
		Where("is_deleted = ?", false).
		Count(&summary.TotalRecords)

	// Category totals
	type categoryRow struct {
		Category string
		Total    float64
		Count    int64
	}
	var rows []categoryRow
	database.DB.Model(&models.FinancialRecord{}).
		Where("is_deleted = ?", false).
		Select("category, SUM(amount) as total, COUNT(*) as count").
		Group("category").
		Order("total DESC").
		Scan(&rows)

	for _, r := range rows {
		summary.CategoryTotals = append(summary.CategoryTotals, models.CategoryTotal{
			Category: r.Category,
			Total:    r.Total,
			Count:    r.Count,
		})
	}

	// Monthly trends (last 12 months)
	type monthRow struct {
		Month   string
		Type    models.TransactionType
		Total   float64
	}
	var monthRows []monthRow
	database.DB.Model(&models.FinancialRecord{}).
		Where("is_deleted = ?", false).
		Select("strftime('%Y-%m', date) as month, type, SUM(amount) as total").
		Group("month, type").
		Order("month DESC").
		Limit(24).
		Scan(&monthRows)

	// Build a map of month -> trend
	trendMap := map[string]*models.MonthlyTrend{}
	for _, r := range monthRows {
		if _, ok := trendMap[r.Month]; !ok {
			trendMap[r.Month] = &models.MonthlyTrend{Month: r.Month}
		}
		if r.Type == models.TypeIncome {
			trendMap[r.Month].Income = r.Total
		} else {
			trendMap[r.Month].Expense = r.Total
		}
	}
	for _, t := range trendMap {
		t.Net = t.Income - t.Expense
		summary.MonthlyTrends = append(summary.MonthlyTrends, *t)
	}

	// Recent activity (last 10 records)
	database.DB.Model(&models.FinancialRecord{}).
		Where("is_deleted = ?", false).
		Preload("User").
		Order("created_at DESC").
		Limit(10).
		Find(&summary.RecentActivity)

	return summary, nil
}

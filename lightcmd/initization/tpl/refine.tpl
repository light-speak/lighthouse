
// ApplyPagination 应用分页
func ApplyPagination(query *gorm.DB, pagination *models.PaginationInput) *gorm.DB {
	if pagination == nil {
		return query
	}

	current := 1
	pageSize := 10

	if pagination.Current != nil && *pagination.Current > 0 {
		current = *pagination.Current
	}
	if pagination.PageSize != nil && *pagination.PageSize > 0 {
		pageSize = *pagination.PageSize
	}

	offset := (current - 1) * pageSize
	return query.Offset(offset).Limit(pageSize)
}

// ApplySorters 应用排序
func ApplySorters(query *gorm.DB, sorters []*models.SorterInput) *gorm.DB {
	if len(sorters) == 0 {
		// 默认按创建时间倒序
		return query.Order("created_at DESC")
	}

	for _, sorter := range sorters {
		// 将驼峰转为下划线
		field := ToSnakeCase(sorter.Field)
		order := "ASC"
		if sorter.Order == models.SortOrderDesc {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", field, order))
	}

	return query
}

// ApplyFilters 应用过滤条件
func ApplyFilters(query *gorm.DB, filters []*models.FilterInput) *gorm.DB {
	if len(filters) == 0 {
		return query
	}

	for _, filter := range filters {
		field := ToSnakeCase(filter.Field)
		value := ""
		if filter.Value != nil {
			value = *filter.Value
		}

		switch filter.Operator {
		case models.FilterOperatorEq:
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		case models.FilterOperatorNe:
			query = query.Where(fmt.Sprintf("%s != ?", field), value)
		case models.FilterOperatorLt:
			query = query.Where(fmt.Sprintf("%s < ?", field), value)
		case models.FilterOperatorLte:
			query = query.Where(fmt.Sprintf("%s <= ?", field), value)
		case models.FilterOperatorGt:
			query = query.Where(fmt.Sprintf("%s > ?", field), value)
		case models.FilterOperatorGte:
			query = query.Where(fmt.Sprintf("%s >= ?", field), value)
		case models.FilterOperatorContains:
			query = query.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value+"%")
		case models.FilterOperatorStartsWith:
			query = query.Where(fmt.Sprintf("%s LIKE ?", field), value+"%")
		case models.FilterOperatorEndsWith:
			query = query.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value)
		case models.FilterOperatorIn:
			// value 是逗号分隔的列表
			values := strings.Split(value, ",")
			query = query.Where(fmt.Sprintf("%s IN ?", field), values)
		case models.FilterOperatorNin:
			values := strings.Split(value, ",")
			query = query.Where(fmt.Sprintf("%s NOT IN ?", field), values)
		case models.FilterOperatorNull:
			query = query.Where(fmt.Sprintf("%s IS NULL", field))
		case models.FilterOperatorNnull:
			query = query.Where(fmt.Sprintf("%s IS NOT NULL", field))
		case models.FilterOperatorBetween:
			// value 格式: "start,end"
			parts := strings.SplitN(value, ",", 2)
			if len(parts) == 2 {
				query = query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), parts[0], parts[1])
			}
		}
	}

	return query
}

// ToSnakeCase 驼峰转下划线
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(r + 32) // 转小写
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// GetPaginationInfo 获取分页信息
func GetPaginationInfo(pagination *models.PaginationInput) (current, pageSize int) {
	current = 1
	pageSize = 10

	if pagination != nil {
		if pagination.Current != nil && *pagination.Current > 0 {
			current = *pagination.Current
		}
		if pagination.PageSize != nil && *pagination.PageSize > 0 {
			pageSize = *pagination.PageSize
		}
	}

	return
}

// ParseFilterValue 解析过滤值为具体类型
func ParseFilterValueInt(value string) int {
	v, _ := strconv.Atoi(value)
	return v
}

func ParseFilterValueInt64(value string) int64 {
	v, _ := strconv.ParseInt(value, 10, 64)
	return v
}

func ParseFilterValueBool(value string) bool {
	return value == "true" || value == "1"
}

func ParseFilterValueTime(value string) time.Time {
	// 尝试多种格式
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, value); err == nil {
			return t
		}
	}
	return time.Time{}
}

// FindFilterValue 从过滤列表中找到指定字段的值
func FindFilterValue(filters []*models.FilterInput, field string) *string {
	for _, f := range filters {
		if f.Field == field {
			return f.Value
		}
	}
	return nil
}

// FindFilterOperator 从过滤列表中找到指定字段的操作符
func FindFilter(filters []*models.FilterInput, field string) *models.FilterInput {
	for _, f := range filters {
		if f.Field == field {
			return f
		}
	}
	return nil
}

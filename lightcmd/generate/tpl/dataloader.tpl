{{- range $model, $fields := .Loader }}
{{- $modelName := ($model | ucFirst) }}
{{- range $field := $fields }}
{{- $loaderName := printf "%s%sLoader" $modelName ($field | ucFirst) }}
{{- $listLoaderName := printf "%s%sListLoader" $modelName ($field | ucFirst) }}

type {{ $loaderName }} struct {
	loader *dataloadgen.Loader[uint, *{{ $modelName }}]
}

type {{ $listLoaderName }} struct {
	loader *dataloadgen.Loader[uint, []*{{ $modelName }}]
}

func (l *{{ $loaderName }}) IsLoader() bool {
	return true
}

func (l *{{ $listLoaderName }}) IsLoader() bool {
	return true
}

func (l *{{ $loaderName }}) Name() string {
	return "{{ $loaderName }}"
}

func (l *{{ $listLoaderName }}) Name() string {
	return "{{ $listLoaderName }}"
}

func (l *{{ $loaderName }}) NewLoader(db *gorm.DB) dataloader.Dataloader {
	return New{{ $loaderName }}(db)
}

func (l *{{ $listLoaderName }}) NewLoader(db *gorm.DB) dataloader.Dataloader {
	return New{{ $listLoaderName }}(db)
}

func New{{ $loaderName }}(db *gorm.DB) *{{ $loaderName }} {
	return &{{ $loaderName }}{
		loader: dataloadgen.NewLoader(func(ctx context.Context, keys []uint) ([]*{{ $modelName }}, []error) {
			return fetch{{ $modelName }}sBy{{ $field | ucFirst }}(ctx, db, keys)
		}),
	}
}

func New{{ $listLoaderName }}(db *gorm.DB) *{{ $listLoaderName }} {
	return &{{ $listLoaderName }}{
		loader: dataloadgen.NewLoader(func(ctx context.Context, keys []uint) ([][]*{{ $modelName }}, []error) {
			return fetch{{ $modelName }}sListBy{{ $field | ucFirst }}(ctx, db, keys)
		}),
	}
}

func Get{{ $loaderName }}(ctx context.Context) *{{ $loaderName }} {
	return dataloader.GetLoaderFromCtx(ctx, "{{ $loaderName }}").(*{{ $loaderName }})
}

func Get{{ $listLoaderName }}(ctx context.Context) *{{ $listLoaderName }} {
	return dataloader.GetLoaderFromCtx(ctx, "{{ $listLoaderName }}").(*{{ $listLoaderName }})
}

func (l *{{ $loaderName }}) Load(ctx context.Context, id uint) (*{{ $modelName }}, error) {
	return l.loader.Load(ctx, id)
}

func (l *{{ $loaderName }}) LoadAll(ctx context.Context, ids []uint) ([]*{{ $modelName }}, error) {
	return l.loader.LoadAll(ctx, ids)
}

func (l *{{ $listLoaderName }}) Load(ctx context.Context, id uint) ([]*{{ $modelName }}, error) {
	return l.loader.Load(ctx, id)
}

func (l *{{ $listLoaderName }}) LoadAll(ctx context.Context, ids []uint) ([][]*{{ $modelName }}, error) {
	return l.loader.LoadAll(ctx, ids)
}

func fetch{{ $modelName }}sBy{{ $field | ucFirst }}(ctx context.Context, db *gorm.DB, keys []uint) ([]*{{ $modelName }}, []error) {
	items := []*{{ $modelName }}{}
	result := db.WithContext(ctx).Where("{{ $field | snakeCase }} IN (?)", keys).Find(&items)
	if result.Error != nil {
		// Return error for each key when query fails
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = result.Error
		}
		return nil, errs
	}

	// Create a map for easier lookup
	itemByKey := map[uint]*{{ $modelName }}{}
	for _, item := range items {
		itemByKey[item.{{ $field | camelCaseWithSpecial }}] = item
	}

	// Return results and errors for each key
	results := make([]*{{ $modelName }}, len(keys))
	errs := make([]error, len(keys))
	for i, key := range keys {
		if item, ok := itemByKey[key]; ok {
			results[i] = item
		} else {
			errs[i] = lighterrors.NewNotFoundError("{{ $modelName }} not found")
		}
	}

	return results, errs
}

func fetch{{ $modelName }}sListBy{{ $field | ucFirst }}(ctx context.Context, db *gorm.DB, keys []uint) ([][]*{{ $modelName }}, []error) {
	items := []*{{ $modelName }}{}
	result := db.WithContext(ctx).Where("{{ $field | snakeCase}} IN (?)", keys).Find(&items)
	if result.Error != nil {
		// Return error for each key when query fails
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = result.Error
		}
		return nil, errs
	}

	// Create a map for easier lookup
	itemsByKey := map[uint][]*{{ $modelName }}{}
	for _, item := range items {
		key := item.{{ $field | camelCaseWithSpecial }}
		itemsByKey[key] = append(itemsByKey[key], item)
	}

	// Return results and errors for each key
	results := make([][]*{{ $modelName }}, len(keys))
	errs := make([]error, len(keys))
	for i, key := range keys {
		if items, ok := itemsByKey[key]; ok {
			results[i] = items
		} else {
			results[i] = []*{{ $modelName }}{} // Return empty slice instead of error for lists
		}
	}

	return results, errs
}

{{- range $extraKey := index $.Extra $model }}
{{- $extraLoaderName := printf "%s%sWith%sLoader" $modelName ($field | ucFirst) ($extraKey | ucFirst) }}
{{- $extraListLoaderName := printf "%s%sWith%sListLoader" $modelName ($field | ucFirst) ($extraKey | ucFirst) }}

type {{ $extraLoaderName }} struct {
	loader *dataloadgen.Loader[struct{ID uint; ExtraID interface{}}, *{{ $modelName }}]
}

type {{ $extraListLoaderName }} struct {
	loader *dataloadgen.Loader[struct{ID uint; ExtraID interface{}}, []*{{ $modelName }}]
}

func (l *{{ $extraLoaderName }}) IsLoader() bool {
	return true
}

func (l *{{ $extraListLoaderName }}) IsLoader() bool {
	return true
}

func (l *{{ $extraLoaderName }}) Name() string {
	return "{{ $extraLoaderName }}"
}

func (l *{{ $extraListLoaderName }}) Name() string {
	return "{{ $extraListLoaderName }}"
}

func (l *{{ $extraLoaderName }}) NewLoader(db *gorm.DB) dataloader.Dataloader {
	return New{{ $extraLoaderName }}(db)
}

func (l *{{ $extraListLoaderName }}) NewLoader(db *gorm.DB) dataloader.Dataloader {
	return New{{ $extraListLoaderName }}(db)
}

func New{{ $extraLoaderName }}(db *gorm.DB) *{{ $extraLoaderName }} {
	return &{{ $extraLoaderName }}{
		loader: dataloadgen.NewLoader(func(ctx context.Context, keys []struct{ID uint; ExtraID interface{}}) ([]*{{ $modelName }}, []error) {
			return fetch{{ $modelName }}sBy{{ $field | ucFirst }}With{{ $extraKey | ucFirst }}(ctx, db, keys)
		}),
	}
}

func New{{ $extraListLoaderName }}(db *gorm.DB) *{{ $extraListLoaderName }} {
	return &{{ $extraListLoaderName }}{
		loader: dataloadgen.NewLoader(func(ctx context.Context, keys []struct{ID uint; ExtraID interface{}}) ([][]*{{ $modelName }}, []error) {
			return fetch{{ $modelName }}sListBy{{ $field | ucFirst }}With{{ $extraKey | ucFirst }}(ctx, db, keys)
		}),
	}
}

func Get{{ $extraLoaderName }}(ctx context.Context) *{{ $extraLoaderName }} {
	return dataloader.GetLoaderFromCtx(ctx, "{{ $extraLoaderName }}").(*{{ $extraLoaderName }})
}

func Get{{ $extraListLoaderName }}(ctx context.Context) *{{ $extraListLoaderName }} {
	return dataloader.GetLoaderFromCtx(ctx, "{{ $extraListLoaderName }}").(*{{ $extraListLoaderName }})
}

func (l *{{ $extraLoaderName }}) Load(ctx context.Context, id uint, extraID interface{}) (*{{ $modelName }}, error) {
	return l.loader.Load(ctx, struct{ID uint; ExtraID interface{}}{ID: id, ExtraID: extraID})
}

func (l *{{ $extraLoaderName }}) LoadAll(ctx context.Context, ids []uint, extraIDs []interface{}) ([]*{{ $modelName }}, error) {
	keys := make([]struct{ID uint; ExtraID interface{}}, len(ids))
	for i := range ids {
		keys[i] = struct{ID uint; ExtraID interface{}}{ID: ids[i], ExtraID: extraIDs[i]}
	}
	return l.loader.LoadAll(ctx, keys)
}

func (l *{{ $extraListLoaderName }}) Load(ctx context.Context, id uint, extraID interface{}) ([]*{{ $modelName }}, error) {
	return l.loader.Load(ctx, struct{ID uint; ExtraID interface{}}{ID: id, ExtraID: extraID})
}

func (l *{{ $extraListLoaderName }}) LoadAll(ctx context.Context, ids []uint, extraIDs []interface{}) ([][]*{{ $modelName }}, error) {
	keys := make([]struct{ID uint; ExtraID interface{}}, len(ids))
	for i := range ids {
		keys[i] = struct{ID uint; ExtraID interface{}}{ID: ids[i], ExtraID: extraIDs[i]}
	}
	return l.loader.LoadAll(ctx, keys)
}

func fetch{{ $modelName }}sBy{{ $field | ucFirst }}With{{ $extraKey | ucFirst }}(ctx context.Context, db *gorm.DB, keys []struct{ID uint; ExtraID interface{}}) ([]*{{ $modelName }}, []error) {
	items := []*{{ $modelName }}{}
	
	// Extract IDs and ExtraIDs
	ids := make([]uint, len(keys))
	extraIDs := make([]interface{}, len(keys))
	for i, key := range keys {
		ids[i] = key.ID
		extraIDs[i] = key.ExtraID
	}
	
	result := db.WithContext(ctx).Where("{{ $field | snakeCase }} IN (?) AND {{ $extraKey | snakeCase }} IN (?)", ids, extraIDs).Find(&items)
	if result.Error != nil {
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = result.Error
		}
		return nil, errs
	}

	// Create a map for easier lookup using composite key
	itemByKey := map[string]*{{ $modelName }}{}
	for _, item := range items {
		key := fmt.Sprintf("%v_%v", item.{{ $field | camelCaseWithSpecial }}, item.{{ $extraKey | camelCaseWithSpecial }})
		itemByKey[key] = item
	}

	results := make([]*{{ $modelName }}, len(keys))
	errs := make([]error, len(keys))
	for i, key := range keys {
		lookupKey := fmt.Sprintf("%v_%v", key.ID, key.ExtraID)
		if item, ok := itemByKey[lookupKey]; ok {
			results[i] = item
		} else {
			errs[i] = lighterrors.NewNotFoundError("{{ $modelName }} not found")
		}
	}

	return results, errs
}

func fetch{{ $modelName }}sListBy{{ $field | ucFirst }}With{{ $extraKey | ucFirst }}(ctx context.Context, db *gorm.DB, keys []struct{ID uint; ExtraID interface{}}) ([][]*{{ $modelName }}, []error) {
	items := []*{{ $modelName }}{}
	
	// Extract IDs and ExtraIDs
	ids := make([]uint, len(keys))
	extraIDs := make([]interface{}, len(keys))
	for i, key := range keys {
		ids[i] = key.ID
		extraIDs[i] = key.ExtraID
	}
	
	result := db.WithContext(ctx).Where("{{ $field | snakeCase }} IN (?) AND {{ $extraKey | snakeCase }} IN (?)", ids, extraIDs).Find(&items)
	if result.Error != nil {
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = result.Error
		}
		return nil, errs
	}

	// Create a map for easier lookup using composite key
	itemsByKey := map[string][]*{{ $modelName }}{}
	for _, item := range items {
		key := fmt.Sprintf("%v_%v", item.{{ $field | camelCaseWithSpecial }}, item.{{ $extraKey | camelCaseWithSpecial }})
		itemsByKey[key] = append(itemsByKey[key], item)
	}

	results := make([][]*{{ $modelName }}, len(keys))
	errs := make([]error, len(keys))
	for i, key := range keys {
		lookupKey := fmt.Sprintf("%v_%v", key.ID, key.ExtraID)
		if items, ok := itemsByKey[lookupKey]; ok {
			results[i] = items
		} else {
			results[i] = []*{{ $modelName }}{} // Return empty slice instead of error for lists
		}
	}

	return results, errs
}

{{- end }}
{{- end }}
{{- end }}

{{- range $model, $morphFields := .MorphTo }}
{{- $modelName := ($model | ucFirst) }}
{{- range $morph := $morphFields }}
{{- range $unionType := $morph.Union }}
{{- $loaderName := printf "%s%sWith%sLoader" $modelName ($morph.Field | ucFirst) $unionType }}
{{- $listLoaderName := printf "%s%sWith%sListLoader" $modelName ($morph.Field | ucFirst) $unionType }}

type {{ $loaderName }} struct {
	loader *dataloadgen.Loader[uint, *{{ $modelName }}]
}

type {{ $listLoaderName }} struct {
	loader *dataloadgen.Loader[uint, []*{{ $modelName }}]
}

func (l *{{ $loaderName }}) IsLoader() bool {
	return true
}

func (l *{{ $listLoaderName }}) IsLoader() bool {
	return true
}

func (l *{{ $loaderName }}) Name() string {
	return "{{ $loaderName }}"
}

func (l *{{ $listLoaderName }}) Name() string {
	return "{{ $listLoaderName }}"
}

func (l *{{ $loaderName }}) NewLoader(db *gorm.DB) dataloader.Dataloader {
	return New{{ $loaderName }}(db)
}

func (l *{{ $listLoaderName }}) NewLoader(db *gorm.DB) dataloader.Dataloader {
	return New{{ $listLoaderName }}(db)
}

func New{{ $loaderName }}(db *gorm.DB) *{{ $loaderName }} {
	return &{{ $loaderName }}{
		loader: dataloadgen.NewLoader(func(ctx context.Context, keys []uint) ([]*{{ $modelName }}, []error) {
			return fetch{{ $modelName }}sBy{{ $morph.Field | ucFirst }}With{{ $unionType }}(ctx, db, keys)
		}),
	}
}

func New{{ $listLoaderName }}(db *gorm.DB) *{{ $listLoaderName }} {
	return &{{ $listLoaderName }}{
		loader: dataloadgen.NewLoader(func(ctx context.Context, keys []uint) ([][]*{{ $modelName }}, []error) {
			return fetch{{ $modelName }}sListBy{{ $morph.Field | ucFirst }}With{{ $unionType }}(ctx, db, keys)
		}),
	}
}

func Get{{ $loaderName }}(ctx context.Context) *{{ $loaderName }} {
	return dataloader.GetLoaderFromCtx(ctx, "{{ $loaderName }}").(*{{ $loaderName }})
}

func Get{{ $listLoaderName }}(ctx context.Context) *{{ $listLoaderName }} {
	return dataloader.GetLoaderFromCtx(ctx, "{{ $listLoaderName }}").(*{{ $listLoaderName }})
}

func (l *{{ $loaderName }}) Load(ctx context.Context, id uint) (*{{ $modelName }}, error) {
	return l.loader.Load(ctx, id)
}

func (l *{{ $loaderName }}) LoadAll(ctx context.Context, ids []uint) ([]*{{ $modelName }}, error) {
	return l.loader.LoadAll(ctx, ids)
}

func (l *{{ $listLoaderName }}) Load(ctx context.Context, id uint) ([]*{{ $modelName }}, error) {
	return l.loader.Load(ctx, id)
}

func (l *{{ $listLoaderName }}) LoadAll(ctx context.Context, ids []uint) ([][]*{{ $modelName }}, error) {
	return l.loader.LoadAll(ctx, ids)
}

func fetch{{ $modelName }}sBy{{ $morph.Field | ucFirst }}With{{ $unionType }}(ctx context.Context, db *gorm.DB, keys []uint) ([]*{{ $modelName }}, []error) {
	items := []*{{ $modelName }}{}
	result := db.WithContext(ctx).Where("{{ $morph.Field | snakeCase }}_id IN (?) AND {{ $morph.Field | snakeCase }}_type = ?", keys, "{{ $unionType | snakeCase }}").Find(&items)
	if result.Error != nil {
		// Return error for each key when query fails
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = result.Error
		}
		return nil, errs
	}

	// Create a map for easier lookup
	itemByKey := map[uint]*{{ $modelName }}{}
	for _, item := range items {
		itemByKey[item.{{ $morph.Field | camelCaseWithSpecial }}ID] = item
	}

	// Return results and errors for each key
	results := make([]*{{ $modelName }}, len(keys))
	errs := make([]error, len(keys))
	for i, key := range keys {
		if item, ok := itemByKey[key]; ok {
			results[i] = item
		} else {
			errs[i] = lighterrors.NewNotFoundError("{{ $modelName }} not found")
		}
	}

	return results, errs
}

func fetch{{ $modelName }}sListBy{{ $morph.Field | ucFirst }}With{{ $unionType }}(ctx context.Context, db *gorm.DB, keys []uint) ([][]*{{ $modelName }}, []error) {
	items := []*{{ $modelName }}{}
	result := db.WithContext(ctx).Where("{{ $morph.Field | snakeCase }}_id IN (?) AND {{ $morph.Field | snakeCase }}_type = ?", keys, "{{ $unionType | snakeCase }}").Find(&items)
	if result.Error != nil {
		// Return error for each key when query fails
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = result.Error
		}
		return nil, errs
	}

	// Create a map for easier lookup
	itemsByKey := map[uint][]*{{ $modelName }}{}
	for _, item := range items {
		key := item.{{ $morph.Field | camelCaseWithSpecial }}ID
		itemsByKey[key] = append(itemsByKey[key], item)
	}

	// Return results and errors for each key
	results := make([][]*{{ $modelName }}, len(keys))
	errs := make([]error, len(keys))
	for i, key := range keys {
		if items, ok := itemsByKey[key]; ok {
			results[i] = items
		} else {
			results[i] = []*{{ $modelName }}{} // Return empty slice instead of error for lists
		}
	}

	return results, errs
}

{{- range $extraKey := index $.Extra $model }}
{{- $extraLoaderName := printf "%s%sWith%sWith%sLoader" $modelName ($morph.Field | ucFirst) $unionType ($extraKey | ucFirst) }}
{{- $extraListLoaderName := printf "%s%sWith%sWith%sListLoader" $modelName ($morph.Field | ucFirst) $unionType ($extraKey | ucFirst) }}

type {{ $extraLoaderName }} struct {
	loader *dataloadgen.Loader[struct{ID uint; ExtraID interface{}}, *{{ $modelName }}]
}

type {{ $extraListLoaderName }} struct {
	loader *dataloadgen.Loader[struct{ID uint; ExtraID interface{}}, []*{{ $modelName }}]
}

func (l *{{ $extraLoaderName }}) IsLoader() bool {
	return true
}

func (l *{{ $extraListLoaderName }}) IsLoader() bool {
	return true
}

func (l *{{ $extraLoaderName }}) Name() string {
	return "{{ $extraLoaderName }}"
}

func (l *{{ $extraListLoaderName }}) Name() string {
	return "{{ $extraListLoaderName }}"
}

func (l *{{ $extraLoaderName }}) NewLoader(db *gorm.DB) dataloader.Dataloader {
	return New{{ $extraLoaderName }}(db)
}

func (l *{{ $extraListLoaderName }}) NewLoader(db *gorm.DB) dataloader.Dataloader {
	return New{{ $extraListLoaderName }}(db)
}

func New{{ $extraLoaderName }}(db *gorm.DB) *{{ $extraLoaderName }} {
	return &{{ $extraLoaderName }}{
		loader: dataloadgen.NewLoader(func(ctx context.Context, keys []struct{ID uint; ExtraID interface{}}) ([]*{{ $modelName }}, []error) {
			return fetch{{ $modelName }}sBy{{ $morph.Field | ucFirst }}With{{ $unionType }}With{{ $extraKey | ucFirst }}(ctx, db, keys)
		}),
	}
}

func New{{ $extraListLoaderName }}(db *gorm.DB) *{{ $extraListLoaderName }} {
	return &{{ $extraListLoaderName }}{
		loader: dataloadgen.NewLoader(func(ctx context.Context, keys []struct{ID uint; ExtraID interface{}}) ([][]*{{ $modelName }}, []error) {
			return fetch{{ $modelName }}sListBy{{ $morph.Field | ucFirst }}With{{ $unionType }}With{{ $extraKey | ucFirst }}(ctx, db, keys)
		}),
	}
}

func Get{{ $extraLoaderName }}(ctx context.Context) *{{ $extraLoaderName }} {
	return dataloader.GetLoaderFromCtx(ctx, "{{ $extraLoaderName }}").(*{{ $extraLoaderName }})
}

func Get{{ $extraListLoaderName }}(ctx context.Context) *{{ $extraListLoaderName }} {
	return dataloader.GetLoaderFromCtx(ctx, "{{ $extraListLoaderName }}").(*{{ $extraListLoaderName }})
}

func (l *{{ $extraLoaderName }}) Load(ctx context.Context, id uint, extraID interface{}) (*{{ $modelName }}, error) {
	return l.loader.Load(ctx, struct{ID uint; ExtraID interface{}}{ID: id, ExtraID: extraID})
}

func (l *{{ $extraLoaderName }}) LoadAll(ctx context.Context, ids []uint, extraIDs []interface{}) ([]*{{ $modelName }}, error) {
	keys := make([]struct{ID uint; ExtraID interface{}}, len(ids))
	for i := range ids {
		keys[i] = struct{ID uint; ExtraID interface{}}{ID: ids[i], ExtraID: extraIDs[i]}
	}
	return l.loader.LoadAll(ctx, keys)
}

func (l *{{ $extraListLoaderName }}) Load(ctx context.Context, id uint, extraID interface{}) ([]*{{ $modelName }}, error) {
	return l.loader.Load(ctx, struct{ID uint; ExtraID interface{}}{ID: id, ExtraID: extraID})
}

func (l *{{ $extraListLoaderName }}) LoadAll(ctx context.Context, ids []uint, extraIDs []interface{}) ([][]*{{ $modelName }}, error) {
	keys := make([]struct{ID uint; ExtraID interface{}}, len(ids))
	for i := range ids {
		keys[i] = struct{ID uint; ExtraID interface{}}{ID: ids[i], ExtraID: extraIDs[i]}
	}
	return l.loader.LoadAll(ctx, keys)
}

func fetch{{ $modelName }}sBy{{ $morph.Field | ucFirst }}With{{ $unionType }}With{{ $extraKey | ucFirst }}(ctx context.Context, db *gorm.DB, keys []struct{ID uint; ExtraID interface{}}) ([]*{{ $modelName }}, []error) {
	items := []*{{ $modelName }}{}
	
	// Extract IDs and ExtraIDs
	ids := make([]uint, len(keys))
	extraIDs := make([]interface{}, len(keys))
	for i, key := range keys {
		ids[i] = key.ID
		extraIDs[i] = key.ExtraID
	}
	
	result := db.WithContext(ctx).Where("{{ $morph.Field | snakeCase }}_id IN (?) AND {{ $morph.Field | snakeCase }}_type = ? AND {{ $extraKey | snakeCase }} IN (?)", ids, "{{ $unionType | snakeCase }}", extraIDs).Find(&items)
	if result.Error != nil {
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = result.Error
		}
		return nil, errs
	}

	// Create a map for easier lookup using composite key
	itemByKey := map[string]*{{ $modelName }}{}
	for _, item := range items {
		key := fmt.Sprintf("%v_%v", item.{{ $morph.Field | camelCaseWithSpecial }}ID, item.{{ $extraKey | camelCaseWithSpecial }})
		itemByKey[key] = item
	}

	results := make([]*{{ $modelName }}, len(keys))
	errs := make([]error, len(keys))
	for i, key := range keys {
		lookupKey := fmt.Sprintf("%v_%v", key.ID, key.ExtraID)
		if item, ok := itemByKey[lookupKey]; ok {
			results[i] = item
		} else {
			errs[i] = lighterrors.NewNotFoundError("{{ $modelName }} not found")
		}
	}

	return results, errs
}

func fetch{{ $modelName }}sListBy{{ $morph.Field | ucFirst }}With{{ $unionType }}With{{ $extraKey | ucFirst }}(ctx context.Context, db *gorm.DB, keys []struct{ID uint; ExtraID interface{}}) ([][]*{{ $modelName }}, []error) {
	items := []*{{ $modelName }}{}
	
	// Extract IDs and ExtraIDs
	ids := make([]uint, len(keys))
	extraIDs := make([]interface{}, len(keys))
	for i, key := range keys {
		ids[i] = key.ID
		extraIDs[i] = key.ExtraID
	}
	
	result := db.WithContext(ctx).Where("{{ $morph.Field | snakeCase }}_id IN (?) AND {{ $morph.Field | snakeCase }}_type = ? AND {{ $extraKey | snakeCase }} IN (?)", ids, "{{ $unionType | snakeCase }}", extraIDs).Find(&items)
	if result.Error != nil {
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = result.Error
		}
		return nil, errs
	}

	// Create a map for easier lookup using composite key
	itemsByKey := map[string][]*{{ $modelName }}{}
	for _, item := range items {
		key := fmt.Sprintf("%v_%v", item.{{ $morph.Field | camelCaseWithSpecial }}ID, item.{{ $extraKey | camelCaseWithSpecial }})
		itemsByKey[key] = append(itemsByKey[key], item)
	}

	results := make([][]*{{ $modelName }}, len(keys))
	errs := make([]error, len(keys))
	for i, key := range keys {
		lookupKey := fmt.Sprintf("%v_%v", key.ID, key.ExtraID)
		if items, ok := itemsByKey[lookupKey]; ok {
			results[i] = items
		} else {
			results[i] = []*{{ $modelName }}{} // Return empty slice instead of error for lists
		}
	}

	return results, errs
}

{{- end }}
{{- end }}
{{- end }}
{{- end }}
func init() {
{{- range $model, $fields := .Loader }}
{{- $modelName := ($model | ucFirst) }}
{{- range $field := $fields }}
	dataloader.RegisterLoader(&{{ $modelName }}{{ $field | ucFirst }}Loader{})
	dataloader.RegisterLoader(&{{ $modelName }}{{ $field | ucFirst }}ListLoader{})
{{- range $extraKey := index $.Extra $model }}
	dataloader.RegisterLoader(&{{ $modelName }}{{ $field | ucFirst }}With{{ $extraKey | ucFirst }}Loader{})
	dataloader.RegisterLoader(&{{ $modelName }}{{ $field | ucFirst }}With{{ $extraKey | ucFirst }}ListLoader{})
{{- end }}
{{- end }}
{{- end }}

{{- range $model, $morphFields := .MorphTo }}
{{- $modelName := ($model | ucFirst) }}
{{- range $morph := $morphFields }}
{{- range $unionType := $morph.Union }}
	dataloader.RegisterLoader(&{{ $modelName }}{{ $morph.Field | ucFirst }}With{{ $unionType }}Loader{})
	dataloader.RegisterLoader(&{{ $modelName }}{{ $morph.Field | ucFirst }}With{{ $unionType }}ListLoader{})
{{- range $extraKey := index $.Extra $model }}
	dataloader.RegisterLoader(&{{ $modelName }}{{ $morph.Field | ucFirst }}With{{ $unionType }}With{{ $extraKey | ucFirst }}Loader{})
	dataloader.RegisterLoader(&{{ $modelName }}{{ $morph.Field | ucFirst }}With{{ $unionType }}With{{ $extraKey | ucFirst }}ListLoader{})
{{- end }}
{{- end }}
{{- end }}
{{- end }}
}

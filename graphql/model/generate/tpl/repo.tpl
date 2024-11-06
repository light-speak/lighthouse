// Generic loader function
func loadEntity[T any](ctx *context.Context, key int64, table string, field string) (*sync.Map, error) {
  data, err := model.GetLoader[int64](model.GetDB(), table, field).Load(key)
  if err != nil {
    return nil, err
  }
  return utils.MapToSyncMap(data), nil
}

// Generic list loader function  
func loadEntityList[T any](ctx *context.Context, key int64, table string, field string) ([]*sync.Map, error) {
  datas, err := model.GetLoader[int64](model.GetDB(), table, field).LoadList(key)
  if err != nil {
    return nil, err
  }
  return utils.MapSliceToSyncMapSlice(datas), nil
}

// Generic query function
func queryEntity[T any](m interface{}, scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return model.GetDB().Model(m).Scopes(scopes...)
}

// Generic first function
func firstEntity[T any](ctx *context.Context, data *sync.Map, enumFieldsFn func(string) func(interface{}) interface{}, 
  model interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error) {
  
  var err error
  var mu sync.Mutex
  
  if data == nil {
    mapData := make(map[string]interface{})
    err = queryEntity[T](model).Scopes(scopes...).First(&mapData).Error
    if err != nil {
      return nil, err
    }
    data = utils.MapToSyncMap(mapData)
  }

  result := &sync.Map{}
  data.Range(func(key, value interface{}) bool {
    k := key.(string)
    if fn := enumFieldsFn(k); fn != nil {
      mu.Lock()
      result.Store(k, fn(value))
      mu.Unlock()
    } else {
      result.Store(k, value)
    }
    return true
  })
  return result, nil
}

// Generic list function
func listEntity[T any](ctx *context.Context, datas []*sync.Map, enumFieldsFn func(string) func(interface{}) interface{}, model interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error) {
  if datas == nil {
    mapDatas := make([]map[string]interface{}, 0)
    err := queryEntity[T](model).Scopes(scopes...).Find(&mapDatas).Error
    if err != nil {
      return nil, err
    }
    datas = utils.MapSliceToSyncMapSlice(mapDatas)
  }

  var mu sync.Mutex
  results := make([]*sync.Map, len(datas))
  
  for i, data := range datas {
    result := &sync.Map{}
    data.Range(func(key, value interface{}) bool {
      k := key.(string)
      if fn := enumFieldsFn(k); fn != nil {
        mu.Lock()
        result.Store(k, fn(value))
        mu.Unlock()
      } else {
        result.Store(k, value)
      }
      return true
    })
    results[i] = result
  }
  
  return results, nil
}

// Generic count function
func countEntity[T any](model interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  var count int64
  err := queryEntity[T](model).Scopes(scopes...).Count(&count).Error
  return count, err
}

{{ range .Nodes }}
{{- $name := .Name -}}
// {{ $name | ucFirst }} functions
func Load__{{ $name | ucFirst }}(ctx *context.Context, key int64, field string) (*sync.Map, error) {
  return loadEntity[models.{{ $name | ucFirst }}](ctx, key, "{{ if ne .Table "" }}{{ .Table }}{{ else }}{{ $name | pluralize | lcFirst }}{{ end }}", field)
}

func LoadList__{{ $name | ucFirst }}(ctx *context.Context, key int64, field string) ([]*sync.Map, error) {
  return loadEntityList[models.{{ $name | ucFirst }}](ctx, key, "{{ if ne .Table "" }}{{ .Table }}{{ else }}{{ $name | pluralize | lcFirst }}{{ end }}", field)
}

func Query__{{ $name | ucFirst }}(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return queryEntity[models.{{ $name | ucFirst }}](&models.{{ $name | ucFirst }}{}, scopes...)
}

func First__{{ $name | ucFirst }}(ctx *context.Context, data *sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error) {
  return firstEntity[models.{{ $name | ucFirst }}](ctx, data, models.{{ $name | ucFirst }}EnumFields, &models.{{ $name | ucFirst }}{}, scopes...)
}

func List__{{ $name | ucFirst }}(ctx *context.Context, datas []*sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error) {
  return listEntity[models.{{ $name | ucFirst }}](ctx, datas, models.{{ $name | ucFirst }}EnumFields, &models.{{ $name | ucFirst }}{}, scopes...)
}

func Count__{{ $name | ucFirst }}(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  return countEntity[models.{{ $name | ucFirst }}](&models.{{ $name | ucFirst }}{}, scopes...)
}
{{ end }}

func init() {
  {{- range .Nodes }}
  model.AddQuickFirst("{{ .Name | ucFirst }}", First__{{ .Name | ucFirst }})
  model.AddQuickList("{{ .Name | ucFirst }}", List__{{ .Name | ucFirst }})
  model.AddQuickLoad("{{ .Name | ucFirst }}", Load__{{ .Name | ucFirst }})
  model.AddQuickLoadList("{{ .Name | ucFirst }}", LoadList__{{ .Name | ucFirst }})
  model.AddQuickCount("{{ .Name | ucFirst }}", Count__{{ .Name | ucFirst }})
  {{- end }}
}

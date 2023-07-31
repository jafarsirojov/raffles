package estate

import (
	"crm/internal/structs"
	"strconv"
)

func (r *repo) getWhereByOptions(options structs.SearchOptionsDTO, offset, limit int) (where string, fields []interface{}) {
	i := 2
	fields = append(fields, offset, limit)

	if options.PriceMin > 0 {
		i++
		where = where + " AND price > $" + strconv.Itoa(i)
		fields = append(fields, options.PriceMin)
	}

	if options.PriceMax > 0 {
		i++
		where = where + " AND price < $" + strconv.Itoa(i)
		fields = append(fields, options.PriceMax)
	}

	if options.BedsMax > 0 {
		i++
		where = where + " AND beds > $" + strconv.Itoa(i)
		fields = append(fields, options.BedsMax)
	}

	if options.BathsMax > 0 {
		i++
		where = where + " AND baths > $" + strconv.Itoa(i)
		fields = append(fields, options.BathsMax)
	}

	if options.PropertyTypes != nil {
		i++
		where = where + " AND property_type = any($" + strconv.Itoa(i) + ")"
		fields = append(fields, options.PropertyTypes)
	}

	if options.SquareFootageMin > 0 {
		i++
		where = where + " AND area_in_meter > $" + strconv.Itoa(i)
		fields = append(fields, options.SquareFootageMin)
	}

	if options.SquareFootageMax > 0 {
		i++
		where = where + " AND area_in_meter < $" + strconv.Itoa(i)
		fields = append(fields, options.SquareFootageMax)
	}

	if options.LotSizeMin > 0 {
		i++
		where = where + " AND lot_size_in_acres > $" + strconv.Itoa(i)
		fields = append(fields, options.LotSizeMin)
	}

	if options.LotSizeMax > 0 {
		i++
		where = where + " AND lot_size_in_acres < $" + strconv.Itoa(i)
		fields = append(fields, options.LotSizeMax)
	}

	if options.YearBuiltMin > 0 {
		i++
		where = where + " AND year_built > $" + strconv.Itoa(i)
		fields = append(fields, options.YearBuiltMin)
	}

	if options.YearBuiltMax > 0 {
		i++
		where = where + " AND year_built < $" + strconv.Itoa(i)
		fields = append(fields, options.YearBuiltMax)
	}

	if options.GarageSpacesMin > 0 {
		i++
		where = where + " AND garage > $" + strconv.Itoa(i)
		fields = append(fields, options.GarageSpacesMin)
	}

	if options.GarageSpacesMax > 0 {
		i++
		where = where + " AND garage < $" + strconv.Itoa(i)
		fields = append(fields, options.GarageSpacesMax)
	}

	if options.Cooling == 1 || options.Cooling == 2 {
		i++
		where = where + " AND cooling = $" + strconv.Itoa(i)
		fields = append(fields, options.Cooling)
	}

	if options.Heating == 1 || options.Heating == 2 {
		i++
		where = where + " AND heating = $" + strconv.Itoa(i)
		fields = append(fields, options.Heating)
	}

	if options.Pool == 1 || options.Pool == 2 {
		i++
		where = where + " AND pool = $" + strconv.Itoa(i)
		fields = append(fields, options.Pool)
	}

	return where + " ORDER BY id DESC offset $1 limit $2 ;", fields
}

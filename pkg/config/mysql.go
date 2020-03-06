package config

import (
	"crud/pkg/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "munzir:munzirdev@tcp(localhost:3306)/maps")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func check_id(input string) bool {
	var rule = regexp.MustCompile("(^[1-9]+$)|(^[1-9]+-[1-9]+$)|(^[1-9]+-[1-9]+-[1-9]+$)|(^[1-9]+-[1-9]+-[1-9]+-[1-9]+$)|(^[1-9]+-[1-9]+-[1-9]+-[1-9]+-[1-9]+$)")
	var ismatch = rule.MatchString(input)
	return ismatch
}

func Mod_Geofence(geo string) string {
	replacer := strings.NewReplacer("[[[", "(((", "]]]", ")))")
	geo = replacer.Replace(geo)
	geo = strings.ReplaceAll(geo, ",", " ")
	geo = strings.ReplaceAll(geo, "] [", ",")
	multipoly_geo := fmt.Sprintf("MULTIPOLYGON%s", geo)
	return multipoly_geo
}

func MapsPost(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var maps models.Maps
	var response models.Response
	var check models.Check
	var status uint16
	w.Header().Set("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&maps)
	input_geo, _ := json.Marshal(maps.Geofence)
	geo := string(input_geo)
	result_geo := Mod_Geofence(geo)
	if maps.Type == "country" {
		fmt.Printf("\nTrying to insert Country data into database\n")
		_, errs := db.Exec("insert into country(name,geofence) values (?,ST_GeomFromText(?))",
			maps.Name,
			result_geo)
		if errs != nil {
			fmt.Println(errs)
		} else {
			fmt.Println("Successfully inserted the Country data into database")
			status = 200
		}
	} else if maps.Type == "province" {
		maps_id, _ := json.Marshal(maps.Id)
		idnya := string(maps_id)
		idnya = idnya[1 : len(idnya)-1]
		res := check_id(idnya)
		if res == true {
			err := db.QueryRow("select country_id from country where country_id = ?",
				idnya).Scan(&check.Id)
			if err != nil {
				fmt.Println(err)
				status = 404
			} else {
				fmt.Printf("\nTrying to insert Province data into database\n")
				_, errs := db.Exec("insert into province(name,geofence,country_id) values (?, ST_GeomFromText(?), ?)",
					maps.Name,
					result_geo,
					idnya)
				if errs != nil {
					fmt.Println(errs)
					status = 500
				} else {
					fmt.Println("Successfully inserted the Province data into database")
					status = 200
				}
			}
		} else {
			status = 400
		}
	} else if maps.Type == "district" {
		maps_id, _ := json.Marshal(maps.Id)
		idnya := string(maps_id)
		idnya = idnya[1 : len(idnya)-1]
		idsplit := strings.Split(idnya, "-")
		res := check_id(idnya)
		if res == true {
			if len(idsplit) == 2 {
				err := db.QueryRow("select c.country_id from country c inner join province p on c.country_id = p.country_id where c.country_id = ? and p.province_id = ?", idsplit[0], idsplit[1]).Scan(&check.Id)
				if err != nil {
					fmt.Println(err)
					status = 404
				} else {
					fmt.Printf("\nTrying to insert District data into database\n")
					_, errs := db.Exec("insert into district(name, geofence, province_id) values (?, ST_GeomFromText(?), ?)",
						maps.Name,
						result_geo,
						idsplit[1])
					if errs != nil {
						fmt.Println(errs)
						status = 500
					} else {
						fmt.Println("Successfully inserted the District data into database")
						status = 200
					}
				}
			} else {
				status = 400
			}
		} else {
			status = 400
		}
	} else if maps.Type == "sub district" {
		maps_id, _ := json.Marshal(maps.Id)
		idnya := string(maps_id)
		idnya = idnya[1 : len(idnya)-1]
		idsplit := strings.Split(idnya, "-")
		res := check_id(idnya)
		if res == true {
			if len(idsplit) == 3 {
				err := db.QueryRow("select c.country_id from country c inner join province p on c.country_id = p.country_id inner join district d on p.province_id = d.province_id where c.country_id = ? and p.province_id = ? and d.district_id = ?", idsplit[0], idsplit[1], idsplit[2]).Scan(&check.Id)
				if err != nil {
					fmt.Println(err)
					status = 404
				} else {
					fmt.Printf("\nTrying to insert Sub District data into database\n")
					_, errs := db.Exec("insert into sub_district(name, geofence, district_id) values (?, ST_GeomFromText(?), ?)",
						maps.Name,
						result_geo,
						idsplit[2])
					if errs != nil {
						fmt.Println(errs)
						status = 500
					} else {
						fmt.Println("Successfully inserted the Sub District data into database")
						status = 200
					}
				}
			} else {
				status = 400
			}
		} else {
			status = 400
		}
	} else if maps.Type == "urban village" {
		maps_id, _ := json.Marshal(maps.Id)
		idnya := string(maps_id)
		idnya = idnya[1 : len(idnya)-1]
		idsplit := strings.Split(idnya, "-")
		res := check_id(idnya)
		if res == true {
			if len(idsplit) == 4 {
				err := db.QueryRow("select c.country_id from country c inner join province p on c.country_id = p.country_id inner join district d on p.province_id = d.province_id inner join sub_district sd on d.district_id = sd.district_id where c.country_id = ? and p.province_id = ? and d.district_id = ? and sd.sdistrict_id = ?",
					idsplit[0],
					idsplit[1],
					idsplit[2],
					idsplit[3]).Scan(&check.Id)
				if err != nil {
					fmt.Println(err)
					status = 404
				} else {
					fmt.Printf("\nTrying to insert Urban Village data into database\n")
					_, errs := db.Exec("insert into urban_village(name, geofence, sdistrict_id) values(?, ST_GeomFromText(?), ?)",
						maps.Name,
						result_geo,
						idsplit[3])
					if errs != nil {
						fmt.Println(errs)
						status = 500
					} else {
						fmt.Println("Successfully inserted the Urban Village data into database")
						status = 200
					}
				}
			} else {
				status = 400
			}
		} else {
			status = 400
		}
	}
	if status == 400 {
		w.WriteHeader(http.StatusBadRequest)
		response.Status = 400
		response.Message = "Bad request. Please input data correctly !"
	} else if status == 200 {
		w.WriteHeader(http.StatusBadRequest)
		response.Status = 200
		response.Message = "Successfully inserted your data !"
	} else if status == 404 {
		w.WriteHeader(http.StatusNotFound)
		response.Status = 404
		response.Message = "Id could not be found"
	} else if status == 500 {
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = 500
		response.Message = "An error occurred"
	}
	json.NewEncoder(w).Encode(response)
}

func MapsGet(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	var status uint16
	var maps models.Maps
	var arr_maps []models.Maps
	var geo_temp models.Temp
	var get models.GetMaps
	var response models.Response_Get
	json.NewDecoder(r.Body).Decode(&get)
	input, _ := json.Marshal(get.Id)
	getid := string(input)
	getid = strings.Trim(getid, "[]")
	perid := strings.Split(getid, ",")
	for i := range perid {
		perid[i] = perid[i][1 : len(perid[i])-1]
		res := check_id(perid[i])
		if res == true {
			strip_split := strings.Split(perid[i], "-")
			if len(strip_split) == 1 {
				var temp_float []float64
				var temp_polygon [][]float64
				var multipolygon [][][]float64
				maps.Type = "country"
				maps.Id = perid[i]
				fmt.Printf("\nGetting the Country data\n")
				rows, err := db.Query("select name, st_astext(geofence) from country where country_id = ?", strip_split[0])
				if err != nil {
					fmt.Println(err)
					status = 500
				} else {
					for rows.Next() {
						if err := rows.Scan(&maps.Name, &geo_temp.Geofence); err != nil {
							fmt.Println(err)
							status = 500
						} else {
							input, _ := json.Marshal(geo_temp.Geofence)
							geo := string(input)
							geo = geo[16 : len(geo)-4]
							polygon := strings.Split(geo, ",")
							for i := range polygon {
								point := strings.Split(polygon[i], " ")
								for ind := range point {
									float_point, _ := strconv.ParseFloat(point[ind], 64)
									temp_float = append(temp_float, float_point)
								}
							}
						}
					}
					for index := 0; index < len(temp_float); index += 2 {
						i := 1
						tempPoint := []float64{temp_float[index], temp_float[i]}
						temp_polygon = append(temp_polygon, tempPoint)
						i += 2
					}
					multipolygon = append(multipolygon, temp_polygon)
					maps.Geofence = multipolygon
					arr_maps = append(arr_maps, maps)
					status = 200
				}
			} else if len(strip_split) == 2 {
				var temp_float []float64
				var temp_polygon [][]float64
				var multipolygon [][][]float64
				maps.Type = "province"
				maps.Id = perid[i]
				fmt.Printf("\nGetting the Province data\n")
				rows, err := db.Query("select p.name,st_astext(p.geofence) from country c inner join province p on c.country_id = p.country_id where c.country_id = ? and p.province_id = ?", strip_split[0], strip_split[1])
				if err != nil {
					fmt.Println(err)
					status = 500
				} else {
					for rows.Next() {
						if err := rows.Scan(&maps.Name, &geo_temp.Geofence); err != nil {
							fmt.Println(err)
							status = 500
						} else {
							input, _ := json.Marshal(geo_temp.Geofence)
							geo := string(input)
							geo = geo[16 : len(geo)-4]
							polygon := strings.Split(geo, ",")
							for i := range polygon {
								point := strings.Split(polygon[i], " ")
								for ind := range point {
									float_point, _ := strconv.ParseFloat(point[ind], 64)
									temp_float = append(temp_float, float_point)
								}
							}
						}
					}
					for index := 0; index < len(temp_float); index += 2 {
						i := 1
						tempPoint := []float64{temp_float[index], temp_float[i]}
						temp_polygon = append(temp_polygon, tempPoint)
						i += 2
					}
					multipolygon = append(multipolygon, temp_polygon)
					maps.Geofence = multipolygon
					arr_maps = append(arr_maps, maps)
					status = 200
				}
			} else if len(strip_split) == 3 {
				var temp_float []float64
				var temp_polygon [][]float64
				var multipolygon [][][]float64
				maps.Type = "district"
				maps.Id = perid[i]
				fmt.Printf("\nGetting the Province data\n")
				rows, err := db.Query("select d.name,st_astext(d.geofence) from country c inner join province p on c.country_id = p.country_id inner join district d on p.province_id = d.province_id where c.country_id = ? and p.province_id = ? and d.district_id = ?", strip_split[0], strip_split[1], strip_split[2])
				if err != nil {
					fmt.Println(err)
					status = 500
				} else {
					for rows.Next() {
						if err := rows.Scan(&maps.Name, &geo_temp.Geofence); err != nil {
							fmt.Println(err)
							status = 500
						} else {
							input, _ := json.Marshal(geo_temp.Geofence)
							geo := string(input)
							geo = geo[16 : len(geo)-4]
							polygon := strings.Split(geo, ",")
							for i := range polygon {
								point := strings.Split(polygon[i], " ")
								for ind := range point {
									float_point, _ := strconv.ParseFloat(point[ind], 64)
									temp_float = append(temp_float, float_point)
								}
							}
						}
					}
					for index := 0; index < len(temp_float); index += 2 {
						i := 1
						tempPoint := []float64{temp_float[index], temp_float[i]}
						temp_polygon = append(temp_polygon, tempPoint)
						i += 2
					}
					multipolygon = append(multipolygon, temp_polygon)
					maps.Geofence = multipolygon
					arr_maps = append(arr_maps, maps)
					status = 200
				}
			} else if len(strip_split) == 4 {
				var temp_float []float64
				var temp_polygon [][]float64
				var multipolygon [][][]float64
				maps.Type = "sub district"
				maps.Id = perid[i]
				fmt.Printf("\nGetting the Province data\n")
				rows, err := db.Query("select sd.name,st_astext(sd.geofence) from country c inner join province p on c.country_id = p.country_id inner join district d on p.province_id = d.province_id inner join sub_district sd on d.district_id = sd.district_id where c.country_id = ? and p.province_id = ? and d.district_id = ? and sd.sdistrict_id = ?", strip_split[0], strip_split[1], strip_split[2], strip_split[3])
				if err != nil {
					fmt.Println(err)
					status = 500
				} else {
					for rows.Next() {
						if err := rows.Scan(&maps.Name, &geo_temp.Geofence); err != nil {
							fmt.Println(err)
							status = 500
						} else {
							input, _ := json.Marshal(geo_temp.Geofence)
							geo := string(input)
							geo = geo[16 : len(geo)-4]
							polygon := strings.Split(geo, ",")
							for i := range polygon {
								point := strings.Split(polygon[i], " ")
								for ind := range point {
									float_point, _ := strconv.ParseFloat(point[ind], 64)
									temp_float = append(temp_float, float_point)
								}
							}
						}
					}
					for index := 0; index < len(temp_float); index += 2 {
						i := 1
						tempPoint := []float64{temp_float[index], temp_float[i]}
						temp_polygon = append(temp_polygon, tempPoint)
						i += 2
					}
					multipolygon = append(multipolygon, temp_polygon)
					maps.Geofence = multipolygon
					arr_maps = append(arr_maps, maps)
					status = 200
				}
			} else if len(strip_split) == 5 {
				var temp_float []float64
				var temp_polygon [][]float64
				var multipolygon [][][]float64
				maps.Type = "urban village"
				maps.Id = perid[i]
				fmt.Printf("\nGetting the Province data\n")
				rows, err := db.Query("select uv.name,st_astext(uv.geofence) from country c inner join province p on c.country_id = p.country_id inner join district d on p.province_id = d.province_id inner join sub_district sd on d.district_id = sd.district_id inner join urban_village uv on sd.sdistrict_id = uv.sdistrict_id where c.country_id = ? and p.province_id = ? and d.district_id = ? and sd.sdistrict_id = ? and uv.uvillage_id = ?", strip_split[0], strip_split[1], strip_split[2], strip_split[3], strip_split[4])
				if err != nil {
					fmt.Println(err)
					status = 500
				} else {
					for rows.Next() {
						if err := rows.Scan(&maps.Name, &geo_temp.Geofence); err != nil {
							fmt.Println(err)
							status = 500
						} else {
							input, _ := json.Marshal(geo_temp.Geofence)
							geo := string(input)
							geo = geo[16 : len(geo)-4]
							polygon := strings.Split(geo, ",")
							for i := range polygon {
								point := strings.Split(polygon[i], " ")
								for ind := range point {
									float_point, _ := strconv.ParseFloat(point[ind], 64)
									temp_float = append(temp_float, float_point)
								}
							}
						}
					}
					for index := 0; index < len(temp_float); index += 2 {
						i := 1
						tempPoint := []float64{temp_float[index], temp_float[i]}
						temp_polygon = append(temp_polygon, tempPoint)
						i += 2
					}
					multipolygon = append(multipolygon, temp_polygon)
					maps.Geofence = multipolygon
					arr_maps = append(arr_maps, maps)
					status = 200
				}
			}
		} else {
			status = 400
		}
	}
	if status == 400 {
		w.WriteHeader(http.StatusBadRequest)
		response.Status = 400
		response.Message = "Bad request. Please input data correctly !"
	} else if status == 200 {
		w.WriteHeader(http.StatusBadRequest)
		response.Status = 200
		response.Message = "Successfully Get your data !"
		response.Data = arr_maps
	} else if status == 404 {
		w.WriteHeader(http.StatusNotFound)
		response.Status = 404
		response.Message = "Id could not be found"
	} else if status == 500 {
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = 500
		response.Message = "An error occurred"
	}
	json.NewEncoder(w).Encode(response)
}

func MapsPut(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var putmaps models.PutMaps
	var response models.Response
	var errmsg bool
	var field string
	var table string
	json.NewDecoder(r.Body).Decode(&putmaps)
	input, _ := json.Marshal(putmaps.Id)
	idnya := string(input)
	idnya = idnya[1 : len(idnya)-1]
	resid := strings.Split(idnya, "-")
	if len(resid) == 1 {
		table = "Country"
		if len(putmaps.Name) > 0 && len(putmaps.Geofence) == 0 {
			fmt.Printf("\nTrying to update the name of the Country data\n")
			_, errs := db.Exec("update country set name=? where country_id = ? ",
				putmaps.Name,
				resid[0])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the name of Country data")
				errmsg = false
				field = "Name"
			}
		} else if len(putmaps.Geofence) > 0 && len(putmaps.Name) == 0 {
			fmt.Printf("\nTrying to update the geofence of the Country data\n")
			input_geo, _ := json.Marshal(putmaps.Geofence)
			geo := string(input_geo)
			result_geo := Mod_Geofence(geo)
			_, errs := db.Exec("update country set geofence=ST_GeomFromText(?) where country_id = ?",
				result_geo,
				resid[0])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the geofence of Country data")
				errmsg = false
				field = "Geofence"
			}
		} else if len(putmaps.Geofence) > 0 && len(putmaps.Name) > 0 {
			fmt.Printf("\nTrying to update the name and geofence of the Country data")
			input_geo, _ := json.Marshal(putmaps.Geofence)
			geo := string(input_geo)
			result_geo := Mod_Geofence(geo)
			_, errs := db.Exec("update country set name = ?, geofence=ST_GeomFromText(?) where country_id = ?",
				putmaps.Name,
				result_geo,
				resid[0])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the name and geofence of Country data")
				errmsg = false
				field = "Name and Geofence"
			}
		} else {
			errmsg = true
		}
	} else if len(resid) == 2 {
		table = "Province"
		if len(putmaps.Name) > 0 && len(putmaps.Geofence) == 0 {
			fmt.Printf("\nTrying to update the name of the Province data\n")
			_, errs := db.Exec("update province set name=? where province_id = ?",
				putmaps.Name,
				resid[0])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the name of Province data")
				errmsg = false
				field = "Name"
			}
		} else if len(putmaps.Geofence) > 0 && len(putmaps.Name) == 0 {
			fmt.Printf("\nTrying to update the geofence of the Province data\n")
			input_geo, _ := json.Marshal(putmaps.Geofence)
			geo := string(input_geo)
			result_geo := Mod_Geofence(geo)
			_, errs := db.Exec("update province set geofence = st_geomfromtext(?) where province_id = ?",
				result_geo,
				resid[1])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the geofence of Province data")
				errmsg = false
				field = "Geofence"
			}
		} else if len(putmaps.Geofence) > 0 && len(putmaps.Name) > 0 {
			fmt.Printf("\nTrying to update the name and geofence of the Province data\n")
			input_geo, _ := json.Marshal(putmaps.Geofence)
			geo := string(input_geo)
			result_geo := Mod_Geofence(geo)
			_, errs := db.Exec("update province set name = ?, geofence=st_geomfromtext(?) where province_id = ?",
				putmaps.Name,
				result_geo,
				resid[1])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the name and geofence of province data")
				errmsg = false
				field = "Name and Geofence"
			}
		} else {
			errmsg = true
		}
	} else if len(resid) == 3 {
		table = "District"
		if len(putmaps.Name) > 0 && len(putmaps.Geofence) == 0 {
			fmt.Printf("\nTrying to update the name of the District data\n")
			_, errs := db.Exec("update district set name = ? where district_id = ?",
				putmaps.Name,
				resid[2])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the name of District data")
				errmsg = false
				field = "Name"
			}
		} else if len(putmaps.Geofence) > 0 && len(putmaps.Name) == 0 {
			fmt.Printf("\nTrying to update the geofence of the District data\n")
			input_geo, _ := json.Marshal(putmaps.Geofence)
			geo := string(input_geo)
			result_geo := Mod_Geofence(geo)
			_, errs := db.Exec("update district set geofence = st_geomfromtext(?) where district_id = ?",
				result_geo,
				resid[2])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the name and geofence data of District data")
				errmsg = false
				field = "Geofence"
			}
		} else if len(putmaps.Geofence) > 0 && len(putmaps.Name) > 0 {
			fmt.Printf("\nTrying to update the name and geofence of District data\n")
			input_geo, _ := json.Marshal(putmaps.Geofence)
			geo := string(input_geo)
			result_geo := Mod_Geofence(geo)
			_, errs := db.Exec("Update district set name = ?, geofence = st_geomfromtext(?) where district_id = ?",
				putmaps.Name,
				result_geo,
				resid[2])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the name and geofence of District data")
				errmsg = false
				field = "Name and Geofence"
			}
		} else {
			errmsg = true
		}
	} else if len(resid) == 4 {
		table = "Sub District"
		if len(putmaps.Name) > 0 && len(putmaps.Geofence) == 0 {
			fmt.Printf("\nTrying to update the name of the Sub District data\n")
			_, errs := db.Exec("update sub_district set name = ? where sdistrict_id = ?",
				putmaps.Name,
				resid[3])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the name of Sub District data")
				errmsg = false
				field = "Name"
			}
		} else if len(putmaps.Geofence) > 0 && len(putmaps.Name) == 0 {
			fmt.Printf("\nTrying to update the geofence of the Sub District data\n")
			input_geo, _ := json.Marshal(putmaps.Geofence)
			geo := string(input_geo)
			result_geo := Mod_Geofence(geo)
			_, errs := db.Exec("update sub_district set geofence = st_geomfromtext(?) where sdistrict_id = ?",
				result_geo,
				resid[3])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the geofence of the Sub District data")
				errmsg = false
				field = "Geofence"
			}
		} else if len(putmaps.Geofence) > 0 && len(putmaps.Name) > 0 {
			fmt.Printf("\nTrying to update the name and geofence of the Sub District data\n")
			input_geo, _ := json.Marshal(putmaps.Geofence)
			geo := string(input_geo)
			result_geo := Mod_Geofence(geo)
			_, errs := db.Exec("update sub_district set name = ?, geofence = st_geomfromtext(?) where sdistrict_id = ?",
				putmaps.Name,
				result_geo,
				resid[3])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the name and geofence of the Sub District data")
				errmsg = false
				field = "Name and Geofence"
			}
		} else {
			errmsg = true
		}
	} else if len(resid) == 5 {
		table = "Urban Village"
		if len(putmaps.Name) > 0 && len(putmaps.Geofence) == 0 {
			fmt.Printf("\nTrying to update the name of the Urban Village data\n")
			_, errs := db.Exec("Update urban_village set name = ? where uvillage_id = ?",
				putmaps.Name,
				resid[4])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the name of Urban Village data")
				errmsg = false
				field = "Name"
			}
		} else if len(putmaps.Geofence) > 0 && len(putmaps.Name) == 0 {
			fmt.Printf("\nTrying to update the geofence of Urban village data\n")
			input_geo, _ := json.Marshal(putmaps.Geofence)
			geo := string(input_geo)
			result_geo := Mod_Geofence(geo)
			_, errs := db.Exec("update urban_village set geofence = st_geomfromtext(?) where uvillage_id = ?",
				result_geo,
				resid[4])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the geofence of urban village data")
				errmsg = false
				field = "Geofence"
			}
		} else if len(putmaps.Geofence) > 0 && len(putmaps.Name) > 0 {
			fmt.Printf("\nTrying to update the name and geofence of the Urban village data\n")
			input_geo, _ := json.Marshal(putmaps.Geofence)
			geo := string(input_geo)
			result_geo := Mod_Geofence(geo)
			_, errs := db.Exec("update urban_village set name = ?, geofence = st_geomfromtext(?) where uvillage_id = ?",
				putmaps.Name,
				result_geo,
				resid[4])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully updated the name and geofence of the Urban Village data")
				errmsg = false
				field = "Name and Geofence"
			}
		} else {
			errmsg = true
		}
	} else {
		errmsg = true
	}
	if errmsg == true {
		response.Status = 0
		response.Message = "Failed update your data. Please contact admiin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully updated the " + field + " of your " + table + " data"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func MapsDelete(w http.ResponseWriter, r *http.Request) {
	var delmaps models.DeleteMaps
	var response models.Response
	var errmsg bool
	var table string
	db := connect()
	defer db.Close()
	json.NewDecoder(r.Body).Decode(&delmaps)
	input_id, _ := json.Marshal(delmaps.Id)
	idnya := string(input_id)
	idnya = idnya[1 : len(idnya)-1]
	id_split := strings.Split(idnya, ",")
	for i := range id_split {
		id_split[i] = id_split[i][1 : len(id_split[i])-1]
		perid := strings.Split(id_split[i], "-")
		if len(perid) == 1 {
			table = "Country"
			fmt.Printf("\nTrying to delete the Country data\n")
			_, errs := db.Exec("delete from country where country_id = ?",
				perid[len(perid)-1])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully deleted the Country data")
			}
		} else if len(perid) == 2 {
			table = "Province"
			fmt.Printf("\nTrying to delete the Province data\n")
			_, errs := db.Exec("delete from province where province_id = ?",
				perid[len(perid)-1])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				errmsg = false
				fmt.Println("Successfully deleted the Province data")
			}
		} else if len(perid) == 3 {
			table = "District"
			fmt.Printf("\nTrying to delete the District data\n")
			_, errs := db.Exec("delete from district where district_id = ?",
				perid[len(perid)-1])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				errmsg = false
				fmt.Println("Successfully deleted the District data")
			}
		} else if len(perid) == 4 {
			table = "Sub District"
			fmt.Printf("\nTrying to delete the Sub District data\n")
			_, errs := db.Exec("delete from sub_district where sdistrict_id = ?",
				perid[len(perid)-1])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully deleted the Sub District data")
				errmsg = false
			}
		} else if len(perid) == 5 {
			table = "Urban Village"
			fmt.Printf("\nTrying to delete Urban Village data\n")
			_, errs := db.Exec("delete from urban_village where uvillage_id = ?",
				perid[len(perid)-1])
			if errs != nil {
				fmt.Println(errs)
				errmsg = true
			} else {
				fmt.Println("Successfully deleted the Urban Village data")
				errmsg = false
			}
		} else {
			errmsg = true
		}
	}
	if errmsg == true {
		response.Status = 0
		response.Message = "An error occured. Please input data correctly or contact admin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully deleted your " + table + " data"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Test(w http.ResponseWriter, r *http.Request) {
	res := check_id("12-12-12a-12")
	if res == true {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

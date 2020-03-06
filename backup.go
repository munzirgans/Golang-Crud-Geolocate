package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang_api_queue/CRUD/pkg/models"
	"log"
	"net/http"
	"strings"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/maps")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Mod_Geofence(geo string) string{
	replacer := strings.NewReplacer("[[[", "(((", "]]]", ")))")
	geo = replacer.Replace(geo)
	geo = strings.ReplaceAll(geo, ",", " ")
	geo = strings.ReplaceAll(geo, "] [", ",")
	multipoly_geo := fmt.Sprintf("MULTIPOLYGON%s", geo)
	return multipoly_geo
}

func PostCountry(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var negara models.Country
	var response models.Response
	json.NewDecoder(r.Body).Decode(&negara)
	input_country, _:= json.Marshal(negara.Name)
	country := string(input_country)
	country = country[1:len(country)-1]
	input_geo, _:= json.Marshal(negara.Geofence_Country)
	geo := string(input_geo)
	result_geo := Mod_Geofence(geo)
	fmt.Printf("\nTrying to insert Country data\n")
	_, errs := db.Exec("insert into country(name, geofence) values (?, ST_GeomFromText(?))",
		country,
		result_geo)
	if errs != nil {
		response.Status = 0
		response.Message = "An error occurred, Please input your data correctly !"
		fmt.Println(errs)
	}else {
		response.Status = 1
		response.Message = "Success insert your data into database"
		fmt.Println("Success insert data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetCountry(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var country models.GetCountry
	var arr_country []models.GetCountry
	var response models.Response_Country
	var logging bool
	fmt.Printf("\nGetting the Country data from database\n")
	rows, err := db.Query("select country_id,name, st_astext(geofence) from country")
	if err != nil {
		log.Print(err)
		logging = true
	}
	for rows.Next() {
		if err := rows.Scan(&country.Id, &country.Name, &country.Geofence); err != nil {
			log.Fatal(err.Error())
			logging = true
		} else {
			arr_country = append(arr_country, country)
		}
	}
	if logging == true {
		response.Status = 0
		response.Message = "An error occured, Please contact admin !"
		response.Data = nil
	}
	if arr_country != nil {
		response.Status = 1
		response.Message = "Successfully got your Country data"
		response.Data = arr_country
		fmt.Println("Successfully got the Country data")
	} else {
		response.Status = 0
		response.Message = "Success got your Country data. But, the Country data is NULL"
		response.Data = nil
		fmt.Println("The Country data is NULL")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteCountry(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var dropmaps models.DeleteMaps
	var response models.Response
	json.NewDecoder(r.Body).Decode(&dropmaps)
	fmt.Printf("\nDeleting the Country data\n")
	_, errs := db.Exec("delete from country where country_id = ?", dropmaps.Id)
	if errs != nil {
		fmt.Println(errs)
		response.Status = 0
		response.Message = "An error occurred while deleting your data. Please contact admin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully deleted your Country data"
		fmt.Println("Successfully deleted the Country data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PutCountry(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var putcountry models.PutCountry
	var response models.Response
	json.NewDecoder(r.Body).Decode(&putcountry)
	input_geo, _:= json.Marshal(putcountry.Geofence)
	geo := string(input_geo)
	result_geo := Mod_Geofence(geo)
	fmt.Printf("\nUpdating the Country data\n")
	_, errs := db.Exec("Update country set name = ?, geofence = ST_GeomFromText(?) where country_id = ?",
		putcountry.Country,
		result_geo,
		putcountry.Country_Id)
	if errs != nil  {
		fmt.Println(errs)
		response.Status = 0
		response.Message = "An error occurred. Please contact admin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully updated your country data"
		fmt.Println("Successfully updated the Country data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PostProvince(w http.ResponseWriter, r *http.Request)  {
	db := connect()
	defer db.Close()
	var province models.Province
	var response models.Response
	json.NewDecoder(r.Body).Decode(&province)
	input_province, _ := json.Marshal(province.Name)
	prov := string(input_province)
	prov = prov[1:len(prov)-1]
	input_geo, _ := json.Marshal(province.Geofence_Province)
	geo := string(input_geo)
	result_geo := Mod_Geofence(geo)
	input_cid, _:= json.Marshal(province.Country_Id)
	cid := string(input_cid)
	fmt.Printf("\nTrying to insert Province data\n")
	_, errs := db.Exec("insert into province(name, geofence, country_id) values (?, ST_GeomFromText(?), ?)",
		prov,
		result_geo,
		cid)
	if errs != nil{
		response.Status = 0
		response.Message ="An Error occurred, Please input your data correctly or contact admin soon !"
		fmt.Println(errs)
	} else {
		response.Status = 1
		response.Message = "Successfully insert your province data into database"
		fmt.Println("Success insert the province data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteProvince(w http.ResponseWriter, r *http.Request)     {
	db := connect()
	defer db.Close()
	var dropmaps models.DeleteMaps
	var response models.Response
	json.NewDecoder(r.Body).Decode(&dropmaps)
	fmt.Printf("\nDeleting the Province data\n")
	_, errs := db.Exec("delete from province where province_id = ?", dropmaps.Id)
	if errs != nil  {
		fmt.Println(errs)
		response.Status = 0
		response.Message = "An error occurred while deleting your province data. Please contact admin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully deleted your Province data"
		fmt.Println("Successfully deleted the Province data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PutProvince(w http.ResponseWriter, r *http.Request)     {
	db := connect()
	defer db.Close()
	var putprovince models.PutProvince
	var response models.Response
	json.NewDecoder(r.Body).Decode(&putprovince)
	input_geo, _:= json.Marshal(putprovince.Geofence)
	geo := string(input_geo)
	result_geo := Mod_Geofence(geo)
	fmt.Printf("\nUpdating the Province data\n")
	_, errs := db.Exec("Update province set country_id = ?, name = ?, geofence = ST_GeomFromText(?) where province_id = ?",
		putprovince.Country_Id,
		putprovince.Province,
		result_geo,
		putprovince.Province_Id)
	if errs != nil {
		fmt.Println(errs)
		response.Status = 0
		response.Message = "An error occurred. Please contact admin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully updated your Province data"
		fmt.Println("Successfully updated the Province data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetProvince(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var province models.GetProvince
	var arr_province []models.GetProvince
	var response models.Response_Province
	var logging bool
	fmt.Printf("\nGetting the Province data from database\n")
	rows, err := db.Query("select province_id, name, st_astext(geofence), country_id from province")
	if err != nil {
		log.Print(err)
		logging = true
	}
	for rows.Next() {
		if err := rows.Scan(&province.Province_Id, &province.Province, &province.Geofence, &province.Country_Id); err != nil {
			log.Fatal(err.Error())
			logging = true
		} else {
			arr_province = append(arr_province, province)
		}
	}
	if logging == true{
		response.Status = 0
		response.Message = "An error occurred. Please contact admin soon !"
		response.Data = nil
	}
	if arr_province != nil {
		response.Status = 1
		response.Message = "Successfully got your Province data"
		response.Data = arr_province
		fmt.Println("Successfully got the Province data")
	} else {
		response.Status = 0
		response.Message= "Success got your Province data. But, the Province data is NULL"
		response.Data = nil
		fmt.Println("The Province data is NULL")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PostDistrict(w http.ResponseWriter, r *http.Request)     {
	db := connect()
	defer db.Close()
	var district models.District
	var response models.Response
	json.NewDecoder(r.Body).Decode(&district)
	input_pid, _ := json.Marshal(district.Province_Id)
	pid := string(input_pid)
	input_district, _ := json.Marshal(district.Name)
	dis := string(input_district)
	dis = dis[1:len(dis)-1]
	input_geo, _:= json.Marshal(district.Geofence_District)
	geo := string(input_geo)
	result_geo := Mod_Geofence(geo)
	fmt.Printf("\nTrying to insert District data into database\n")
	_, errs := db.Exec("insert into district(name,geofence,province_id) values(?,ST_GeomFromText(?), ?)",
		dis,
		result_geo,
		pid)
	if errs != nil{
		fmt.Println(errs)
		response.Status = 0
		response.Message = "An Error occurred. Please input your data correctly or contact admin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully insert your District data into database"
		fmt.Println("Success insert the District data")
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteDistrict(w http.ResponseWriter, r *http.Request)      {
	db := connect()
	defer db.Close()
	var dropmaps models.DeleteMaps
	var response models.Response
	json.NewDecoder(r.Body).Decode(&dropmaps)
	fmt.Printf("\nDeleting the District data\n")
	_, errs := db.Exec("delete from district where district_id = ?", dropmaps.Id)
	if errs != nil  {
		fmt.Println(errs)
		response.Status = 0
		response.Message = "An error occured while deleting your District data. Please contact admin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully deleted your District data"
		fmt.Println("Successfully deleted the District data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PutDistrict(w http.ResponseWriter, r *http.Request)      {
	db := connect()
	defer db.Close()
	var putdistrict models.PutDistrict
	var response models.Response
	json.NewDecoder(r.Body).Decode(&putdistrict)
	input_geo, _ := json.Marshal(putdistrict.Geofence)
	geo := string(input_geo)
	result_geo := Mod_Geofence(geo)
	fmt.Printf("\nUpdating the District data\n")
	_, errs := db.Exec("Update district set province_id = ?, name = ?, geofence = ST_GeomFromText(?) where district_id = ?",
		putdistrict.Province_Id,
		putdistrict.District,
		result_geo,
		putdistrict.District_Id)
	if errs != nil  {
		fmt.Println(errs)
		response.Status = 0
		response.Message = "An error occurred. Please contact admin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully updated your District data"
		fmt.Println("Successfully updated the District data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetDistrict(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var district models.GetDistrict
	var arr_district []models.GetDistrict
	var response models.Response_District
	var logging bool
	fmt.Printf("\nGetting the District data from database\n")
	rows, err := db.Query("select district_id, name, st_astext(geofence), province_id from district")
	if err != nil  {
		log.Print(err)
		logging = true
	}
	for rows.Next()  {
		if err := rows.Scan(&district.District_Id, &district.District, &district.Geofence, &district.Province_Id); err != nil  {
			log.Fatal(err.Error())
			logging = true
		} else {
			arr_district = append(arr_district, district)
		}
	}
	if logging == true {
		response.Status = 0
		response.Message = "An error occurred. Please contact admin soon !"
		response.Data = nil
	}
	if arr_district != nil{
		response.Status = 1
		response.Message = "Successfully got your District data"
		response.Data = arr_district
		fmt.Println("Successfully got the District data")
	} else {
		response.Status = 0
		response.Message = "Success got your District data. But, the District data is NULL"
		response.Data = nil
		fmt.Println("The District data is NULL")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PostSubDistrict(w http.ResponseWriter, r *http.Request)      {
	db := connect()
	defer db.Close()
	var sub_district models.Sub_District
	var response models.Response
	json.NewDecoder(r.Body).Decode(&sub_district)
	input_did, _:= json.Marshal(sub_district.District_Id)
	did := string(input_did)
	input_sdis, _:= json.Marshal(sub_district.Name)
	sdis := string(input_sdis)
	sdis = sdis[1:len(sdis)-1]
	input_geo, _:= json.Marshal(sub_district.Geofence_SubDistrict)
	geo := string(input_geo)
	result_geo := Mod_Geofence(geo)
	fmt.Printf("\nTrying to insert Sub District data into database")
	_, errs := db.Exec("insert into sub_district(name,geofence,district_id) values (?, ST_GeomFromText(?), ?)",
		sdis,
		result_geo,
		did)
	if errs != nil {
		fmt.Println(errs)
		response.Status = 0
		response.Message = "An error occurred. Please contact admin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully insert your Sub District data into database"
		fmt.Println("Successfully inserted the Sub District data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteSubDistrict(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var dropmaps models.DeleteMaps
	var response models.Response
	json.NewDecoder(r.Body).Decode(&dropmaps)
	fmt.Println("Deleting the Sub District data")
	_, errs := db.Exec("delete from sub_district where sdistrict_id = ?",
		dropmaps.Id)
	if errs != nil  {
		fmt.Println(errs)
		response.Status = 0
		response.Message = "An error occurred while deleting your Sub DIstrict data. Please contact admin soon !"
		fmt.Println("Successfully deleted the Sub District data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PutSubDistrict(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var putsdistrict models.PutSub_District
	var response models.Response
	json.NewDecoder(r.Body).Decode(&putsdistrict)
	input_geo, _ := json.Marshal(putsdistrict.Geofence)
	geo := string(input_geo)
	result_geo := Mod_Geofence(geo)
	fmt.Printf("\nUpdating the Sub District data\n")
	_, errs := db.Exec("Update sub_district set name=?, geofence=ST_GeomFromText(?), district_id = ? where sdistrict_id = ?",
		putsdistrict.Sub_District,
		result_geo,
		putsdistrict.District_Id,
		putsdistrict.Sub_District_Id)
	if errs != nil  {
		fmt.Println(errs)
		response.Status = 0 
		response.Message = "An error occurred. Please contact admin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully updated your Sub District data"
		fmt.Println("Successfully updated the Sub District data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetSubDistrict(w http.ResponseWriter, r *http.Request)      {
	db := connect()
	defer db.Close()
	var sub_district models.GetSubDistrict
	var arr_sdistrict []models.GetSubDistrict
	var response models.Response_Sub_District
	var logging bool
	fmt.Printf("\nGetting the Sub District data from database\n")
	rows, err := db.Query("select sdistrict_id, name, st_astext(geofence),district_id from sub_district")
	if err != nil  {
		log.Print(err)
		logging = true
	}
	for rows.Next()  {
		if err := rows.Scan(&sub_district.Sub_District_Id, &sub_district.Sub_District, &sub_district.Geofence, &sub_district.District_Id); err != nil  {
			log.Fatal(err.Error())
			logging = true
		} else {
			arr_sdistrict = append(arr_sdistrict, sub_district)
		}
	}
	if logging == true {
		response.Status = 0
		response.Message = "An error occurred. Please contact admin soon !"
		response.Data = nil
	}
	if arr_sdistrict != nil{
		response.Status = 1
		response.Message= "Successfully got your Sub District data"
		response.Data = arr_sdistrict
		fmt.Println("Successfully got the Sub District data")
	} else {
		response.Status = 0
		response.Message = "Success got your Sub District data. But, the Sub District data is NULL"
		response.Data = nil
		fmt.Println("The Sub District data is NULL")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PostUrbanVillage(w http.ResponseWriter, r *http.Request) {
	db:= connect()
	defer db.Close()
	var uvillage models.Urban_Village
	var response models.Response
	json.NewDecoder(r.Body).Decode(&uvillage)
	input_sdid, _:= json.Marshal(uvillage.Sub_District_Id)
	sdid := string(input_sdid)
	input_uvi, _ := json.Marshal(uvillage.Name)
	uvi := string(input_uvi)
	uvi = uvi[1:len(uvi)-1]
	input_geo, _:= json.Marshal(uvillage.Geofence_UrbanVillage)
	geo := string(input_geo)
	result_geo := Mod_Geofence(geo)
	fmt.Printf("\nTrying to insert Urban Village data into database\n")
	_, errs := db.Exec("insert into urban_village(name, geofence, sdistrict_id) values (?, ST_GeomFromText(?), ?)",
		uvi,
		result_geo,
		sdid)
	if errs != nil{
		fmt.Println(errs)
		response.Status = 0
		response.Message = "An error occurred. Please contact admin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully inserted the Urban Village data into database"
		fmt.Println("Successfully inserted the Urban village data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteUrbanVillage(w http.ResponseWriter, r *http.Request)      {
	db := connect()
	defer db.Close()
	var dropmaps models.DeleteMaps
	var response models.Response
	json.NewDecoder(r.Body).Decode(&dropmaps)
	fmt.Println("Deleting the urban village data")
	_, errs := db.Exec("delete from urban_village where uvillage_id = ?",
		dropmaps.Id)
	if errs != nil {
		fmt.Println(errs)
		response.Status = 0 
		response.Message = "An error occurred while deleting your Urban Village data. Please contact admin soon !"
		fmt.Println("Successfully deleted the Urban Village data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PutUrbanVillage(w http.ResponseWriter, r *http.Request)      {
	db := connect()
	defer db.Close()
	var putuvillage models.PutUrban_Village
	var response models.Response
	json.NewDecoder(r.Body).Decode(&putuvillage)
	input_geo, _ := json.Marshal(putuvillage.Geofence)
	geo := string(input_geo)
	result_geo := Mod_Geofence(geo)
	fmt.Printf("\nUpdating the Urban Village data\n")
	_, errs := db.Exec("update urban_village set name = ?, geofence = ST_GeomFromText(?), sdistrict_id = ? where uvillage_id = ?",
		putuvillage.Urban_Village,
		result_geo,
		putuvillage.Sub_District_Id,
		putuvillage.Urban_Village_Id)
	if errs != nil {
		fmt.Println(errs)
		response.Status = 0
		response.Message = "An error occurred. Please contact admin soon !"
	} else {
		response.Status = 1
		response.Message = "Successfully updated your Urban Village data"
		fmt.Println("Successfully updated the Urban Village data")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUrbanVillage(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var uvillage models.GetUrbanVillage
	var arr_uvillage []models.GetUrbanVillage
	var response models.Response_Urban_Village
	var logging bool
	fmt.Printf("\nGetting the Urban Village data from database")
	rows, err := db.Query("select uvillage_id, name, st_astext(geofence), sdistrict_id from urban_village")
	if err != nil  {
		log.Print(err)
		logging = true
	}
	for rows.Next()  {
		if err := rows.Scan(&uvillage.Urban_Village_Id, &uvillage.Urban_Village, &uvillage.Geofence, &uvillage.Sub_District_Id); err != nil {
			log.Fatal(err.Error())
			logging = true
		} else {
			arr_uvillage = append(arr_uvillage, uvillage)
		}
	}
	if logging == true {
		response.Status = 0
		response.Message = "An error occurred. Please contact admin soon !"
		response.Data = nil
	}
	if arr_uvillage != nil {
		response.Status = 1
		response.Message = "Successfully got your Urban Village data"
		response.Data = arr_uvillage
		fmt.Println("Successfully got the Urban Village data")
	} else {
		response.Status = 0
		response.Message = "Success got your Urban Village data. But, the Urban Village data is NULL"
		response.Data = nil
		fmt.Println("The Urban Village data is NULL")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// func Testing(w http.ResponseWriter, r *http.Request){
// 	var geofence models.Testing
// 	json.NewDecoder(r.Body).Decode(&geofence)
// 	input, _:= json.Marshal(geofence.Geofence)
// 	str := string(input)
// 	replacer := strings.NewReplacer("[[[", "(((", "]]]", ")))")
// 	str = replacer.Replace(str)
// 	str = strings.ReplaceAll(str, ","," ")
// 	str = strings.ReplaceAll(str, "] [", ",")
// 	geo_format := fmt.Sprintf("MULTIPOLYGON%s", str)
// 	fmt.Println(geo_format)
// }
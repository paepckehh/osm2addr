example output:

```Shell 

mkdir -p data/validated-preload
curl --output data/germany-latest.osm.pbf https://download.geofabrik.de/europe/germany-latest.osm.pbf 
curl --output data/validated-preload/DE.csv https://downloads.suche-postleitzahl.org/v2/public/zuordnung_plz_ort.csv 
go run paepcke.de/osm2addr/cmd/osm2addr@latest DE data/germany-latest.osm.pbf

OSM:Startup               # 2025-03-10 07:44:58.79476069 +0000 UTC m=+0.000533658
OSM:TargetCountry         # DE
OSM:WorkerScale           # 1
OSM:File                  # data/germany-latest.osm.pbf
----------------------------------------------------------------------------------
OSM:PreLoadFile           # data/validated-preload/DE.csv
OSM:PreLoadFile:Done      # 12853
----------------------------------------------------------------------------------
OSM:PBF:File:URL          # https://download.geofabrik.de/europe/germany-updates
OSM:PBF:File:Repl:USM     # 4330
OSM:PBF:File:Repl:TS      # 2025-02-13 21:21:14 +0000 UTC
----------------------------------------------------------------------------------
OSM:PBF:Parsed:Objects    # 411.113.874
OSM:PBF:Parsed:Tags       #  71.720.811
OSM:PBF:Parsed:Country    #   2.792.983
OSM:PBF:Parsed:Street     #   3.905.365
OSM:PBF:Parsed:City       #   3.659.945
OSM:PBF:Parsed:Postcode   #   3.513.835
----------------------------------------------------------------------------------
OSM:PBF:Complete:AddrTags #   2.711.889
OSM:PBF:Uniq:Country      #          10
OSM:PBF:Err:Uniform       #      21.851
OSM:PBF:Err:Country       #           0
OSM:PBF:Err:Postcode      #           0
OSM:PBF:Err:City          #          17
OSM:PBF:Err:Street        #         162
----------------------------------------------------------------------------------
OSM:Corrected:Auto:Cases  #          81
OSM:Corrected:Auto:Total  #      21.132
OSM:Collect:Places:Total  #     278.379
----------------------------------------------------------------------------------
OSM:Writer:JSON           # json/DE/id.json
OSM:Writer:JSON           # json/DE/corrected.json
OSM:Time:Total            # 45.257262208s

```


function changeColumn(x){
    if(x === 'dateTime'){   
        var d = new Date();
        d.setFullYear(d.getFullYear() + 1)
        document.cookie = "timestamp=DateTime; path=/; expires=" + d.toUTCString() + ";";

        ageTitle = document.getElementById("age");
        ageTitle.style.display = "none";
        ageValues = document.getElementsByClassName("ageData");
        for (var i = 0; i < ageValues.length; i++) {
            ageValues[i].style.display = 'none';
        }

        dateTimeTitle = document.getElementById("dateTime");
        dateTimeTitle.style.display  = "table-cell"
        dateTimeValues = document.getElementsByClassName("dateTimeData");
        for (var i = 0; i < dateTimeValues.length; i++) {
            dateTimeValues[i].style.display  = 'table-cell';
        }       
    }
    else {
        var d = new Date();
        d.setFullYear(d.getFullYear() + 1)
        document.cookie = "timestamp=Age; path=/; expires=" + d.toUTCString() + ";";

        dateTimeTitle = document.getElementById("dateTime");
        dateTimeTitle.style.display = "none";
        dateTimeValues = document.getElementsByClassName("dateTimeData");
        for (var i = 0; i < dateTimeValues.length; i++) {
            dateTimeValues[i].style.display = 'none';
        }

        ageTitle = document.getElementById("age");
        ageTitle.style.display = "table-cell"
        ageValues = document.getElementsByClassName("ageData");
        for (var i = 0; i < ageValues.length; i++) {
            ageValues[i].style.display = 'table-cell';
        }
    }
}
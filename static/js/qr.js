qrkod = document.querySelector("kodqr")
pokazkodqr = false
function qr(){
    pokazkodqr = !pokazkodqr;
    if (pokazkodqr){
        qrkod.style.display = "block";
    } else {
        qrkod.style.display = "none";
    }
}
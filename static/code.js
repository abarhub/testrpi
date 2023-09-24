console.log('test1');
console.log('test2');


function envoyer() {
    console.log('test3');

    var checkedValue = document.querySelector('input[name="action"]:checked').value;
    console.log('valeur', checkedValue);
    if (checkedValue) {
        let param = '';
        if (checkedValue === 'minuteur') {
            const timeControl = document.querySelector('input[type="time"]');
            param = timeControl.value;
        }
        fetch("/api/action/" + checkedValue + ((param !== '') ? '?time=' + param : ''))
            .then(response => {
                console.log('response:', response)
            })
            //.then(data => {
                //console.log(data.count)
                //console.log(data.products)
            //})
            .catch(error => console.error(error))
    }
}

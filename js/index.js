bullshit = {
    init: function() {
        document.getElementById("btn-get-bullshit")
            .addEventListener("click", this.handleGetBullshit);
    },

    openButtonLoading: function(elem) {
        elem.classList.add("disabled");
        elem.setAttribute("origin-text", elem.innerText);
        elem.innerText = elem.getAttribute("loading-text");
    },

    closeButtonLoading: function(elem) {
        elem.classList.remove("disabled");
        elem.setAttribute("loadign-text", elem.innerText);
        elem.innerText = elem.getAttribute("origin-text");
    },

    handleGetBullshit: function(e) {
        const topic = document.getElementById("topic").value;
        const minLen = +document.getElementById("minlen").value;
        if (topic === "" || minLen === 0 || minLen % 1 !== 0) return;
        bullshit.openButtonLoading(e.target);

        setTimeout(() => {
            $.ajax({
                url: "https://api.howtobullshit.me/bullshit",
                type: "post",
                contentType: "application/json; charset=utf-8",
                data: JSON.stringify({ Topic: topic, MinLen: minLen }),
                dataType: "text",
                success: function(r) {
                    $("#content").html(r)
                    bullshit.closeButtonLoading(e.target)
                }
            })
        }, 1500);
    }
}

bullshit.init()
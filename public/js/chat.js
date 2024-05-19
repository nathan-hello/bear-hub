
let lockToBottom = true
let messageContainerDiv
let scrollBottomButton
let scrollBottomSpan
let inputElement

document.addEventListener("DOMContentLoaded", function(_) {

        messageContainerDiv = document.getElementById("message-container")
        scrollBottomButton = document.getElementById("scroll-bottom-button")
        scrollBottomSpan = document.getElementById("scroll-bottom-span")
        inputElement = document.getElementById("msg-input-id");

        document.addEventListener("htmx:wsAfterSend", function(_) {
                inputElement.value = "";
        })

        document.addEventListener("htmx:oobAfterSwap", function(event) {
                if (event.detail.target.id === "messages") {

                        messageContainerDiv.addEventListener("wheel", function(_) {
                                if (messageContainerDiv.scrollHeight === undefined) {
                                        return
                                }
                                const { offsetHeight, scrollHeight, scrollTop } = messageContainerDiv
                                if ((scrollHeight <= scrollTop + offsetHeight + 100)) {
                                        controlScroll(true)
                                } else {
                                        controlScroll(false)
                                }
                        })
                        if (lockToBottom) {
                                controlScroll(true)
                        }
                }
        })
})
function controlScroll(stick) {
        if (stick === true) {
                messageContainerDiv.scrollTop = messageContainerDiv.scrollHeight
                scrollBottomButton.setAttribute("hidden", "true")
                scrollBottomSpan.setAttribute("hidden", "true")
                lockToBottom = true
        }
        if (stick === false) {
                scrollBottomButton.removeAttribute("hidden")
                scrollBottomSpan.removeAttribute("hidden")
                lockToBottom = false
        }

}

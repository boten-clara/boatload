package frontend

type IFrontendRetriever interface {
	getFrontend() ([]byte, error)
}

type InMemoryHtmlRetriever struct{}

func (u InMemoryHtmlRetriever) getFrontend() ([]byte, error) {
	return []byte(uploadPage), nil
}

const uploadPage = `<!DOCTYPE html>
<html lang="en" class="no-js">
<head>
  <meta charset="UTF-8">
  <title>Boatload oceanographic meteorological data uploader</title>

  <script type="text/javascript">
	function upload() {
	  let statusElement = document.getElementById("status");
	  let file = document.getElementById("file").files[0]
	  let req = new XMLHttpRequest();
	  let formData = new FormData();

	  formData.append("file", file);                                
      req.open("POST", '/api/upload');
      req.send(formData);
      
	  statusElement.innerText = "sending..."
	  
      req.onreadystatechange = function() {
  		  if (req.readyState === 4) {
			  if (req.status === 200) {
	  				statusElement.innerText = "success"
			  } else {
	  				statusElement.innerText = "error: " + req.responseText
			  }
    	  }
  	  }
	}
  </script>

</head>

<body>
<h1>boatload</h1>
<h2>oceanographic research data uploader v1.0</h2>
<input id="file" type="file" />
<input type = "button" onclick = "upload()" value = "upload">  
<p id="status"></p>
</body>
</html>`

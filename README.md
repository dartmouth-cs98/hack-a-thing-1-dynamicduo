<h1>
GIF a little bit
</h1>
<h3> Divya Kalidindi and Drew Waterman </h3>

<h2>
About:   
</h2>
We built a simple website that allows a user to choose exactly three images and to create a GIF out of them. The images can be of any type or size; our program will scale each of the images so that the resulting GIF's width will be the same as the width of the widest image chosen by the user, and the resulting GIF's height will be the height of the tallest image uploaded by the user. Upon clicking 'Create GIF', the user is taken to a different page, where they will see the resulting GIF that was created from the three images. A copy of each of the three original images also gets saved to the program directory, since that's the only way we could figure out how to do this. Once the user clicks 'Create GIF', they're gonna need to wait a while before anything happens, because our program is just slow. Also, as of know, our program only takes image files in the format jpeg, png, and gif, but images are not restricted by size. 


<h2>
Technologies used/What we learned:  
</h2>
We decided to play around with the Go programming language (also known as golang). We did 0 research about this programming language beforehand, and we decided we wanted to create a simple website with it. Big mistake. We had a lot of trouble trying to use Go for front-end development, and in our opinion, Go is much better suited for server-side applications. Therefore, we learned that it is best to do solid research on a programming language before deciding to do a project with it. This language was new to both of us, so we also got to learn some aspects of an unfamiliar programming language, and we got a little more comfortable with it by the end of this project.

<h2>
To run our program:  
</h2>

Download Go:  
https://golang.org/doc/install?download=go1.9.2.windows-amd64.msi   
We also used a package: from the project directory, enter the command: <code>go get -u github.com/disintegration/imaging  </code>
Once everything has been downloaded:  
<code>
cd src  
</code>  
<code>
go run app.go
</code>  
Then open your browser to `localhost:3000`  
Don't forget that nothing will happen for a while after you click 'Create GIF'. That's fine! Patience is a virtue. Just keep waiting until you're redirected.  

<h2>
Work Distribution:
</h2>

The two of us pair programmed throughout most of this project. We mainly worked on Drew's computer, but both spent time together working on solutions, researching tutorials, looking through documentation, and debugging. Divya styled the webpage, Drew worked on the README, and we coded the solution together.  

<h2>
Workflow/References:  
</h2>

We started off by creating GIFs from static images stored in our project directory. We used code from the following example:  
http://tech.nitoyon.com/en/blog/2016/01/07/go-animated-gif-gen/

We then created a static site, using the following tutorial:  
http://www.alexedwards.net/blog/serving-static-sites-with-go  

This article helped us provide functionality for allowing users to upload files through the website:  
https://astaxie.gitbooks.io/build-web-application-with-golang/en/04.5.html  

The following link helped us to display the GIFs we created on a new webpage:  
https://www.sanarias.com/blog/1214PlayingwithimagesinHTTPresponseingolang, and we got help from another article as well to learn how to encode and decode images properly: https://www.devdungeon.com/content/working-images-go  

For scaling and saving a GIF, we drew from code from the following repository:  
https://github.com/srinathh/goanigiffy/blob/master/goanigiffy.go

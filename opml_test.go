package opml

import (
	"bytes"
	"io"
	"net/url"
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"golang.org/x/net/html/charset"
)

// http://dev.opml.org/examples/specification.opml
var specification = &OPML{
	Version:         "1.0",
	Title:           "specification.xml",
	DateCreated:     parseTime("Thu, 27 Jul 2000 01:20:06 GMT"),
	DateModified:    parseTime("Fri, 15 Sep 2000 09:04:03 GMT"),
	OwnerName:       "Dave Winer",
	OwnerEmail:      "dave@userland.com",
	ExpansionState:  []int{},
	VertScrollState: 1,
	WindowTop:       146,
	WindowLeft:      107,
	WindowBottom:    468,
	WindowRight:     560,
	Outlines: []*Outline{
		{
			Text: "It's XML, of course",
			Outlines: []*Outline{
				{Text: "This page documents the file formats used by Radio UserLand."},
				{Text: "There are two formats, outlineDocument and songList."},
				{Text: "There's a simple XML-RPC interface that allows a user to register with an aggregator. "},
				{Text: "All formats are open and public and may be used for any purpose whatsoever."},
			},
		},
		{
			Text: "outlineDocument",
			Outlines: []*Outline{
				{Text: "All playlists are outlineDocuments. This is the main file format for Radio UserLand. When you create a new file it's saved as an outlineDocument. Because users can save them into the www folder, they can be accessed over the Web, either from a script or a Web browser such as MSIE. (Of course they could be served by any HTTP server, not just the Radio UserLand server.)"},
				{Text: `The outlineDocument format is <a href="http://backend.userland.com/stories/storyReader$53">documented</a> on backend.userland.com. There will no doubt be changes and refinements to the format. One area that needs work is the format for the data attribute on a headline. Currently there are bugs in the way Radio UserLand uses this attribute. (Every headline gets a data attribute, whether or not it links to a song. We need to XMLize this and fit in data not as an attribute but as a legal sub-item. Shouldn't be hard to do, and with this caveat, breakage should be expected.)`},
				{Text: "Radio UserLand can be used to write any kind of document, not just a music playlist. Outlines are great for all kinds of structured documents, specifications, legal briefs, product plans, presentations and stories."},
				{Text: "Several examples of outlineDocuments created with Radio UserLand: play list, specification, presentation."},
			},
		},
		{
			Text: "songList",
			Outlines: []*Outline{
				{Text: `As you're listening to music, Radio UserLand keeps track of what you listen to. Here's a <a href="http://static.userland.com/images/radiodiscuss/userPlaylistSongs.gif">screen shot</a> of the table, user.playlist.songs, that keeps track of the stuff. `},
				{Text: "ctPlays is the number of times the song has been played. ctSeconds is the duration of the song, determined by a heuristic that's pretty accurate. f is the file that contains the MP3, on the local file system. whenFirstPlayed is the time/date the song was played for the first time, whenLastPlayed is the most recent time/date. whenLoaded is when Radio UserLand discovered the file in your MP3 folder."},
				{Text: "Every hour on the hour Radio UserLand generates an XMLization of this table and places it in the userland folder of your www folder, making it available over the Web. (There's no way to turn this feature off, there should be.)"},
				{Text: `Here's an <a href="http://static.userland.com/gems/radiodiscuss/songs.xml">example</a> of the XML file. The mapping between the table and the XMLization should be fairly clear.`},
			},
		},
		{
			Text: "Rules of the road",
			Outlines: []*Outline{
				{Text: "Rules of the road will be determined later, since many of these files will be on users' machines, we want to provide guidelines for bots, aggregators and content systems; and whatever other kinds of applications people think of. Feel free to use the discussion group here to raise issues. "},
			},
		},
	},
}

// http://dev.opml.org/examples/presentation.opml
var presentation = &OPML{
	Version:         "1.0",
	Title:           "presentation.xml",
	DateCreated:     parseTime("Thu, 27 Jul 2000 01:35:52 GMT"),
	DateModified:    parseTime("Fri, 15 Sep 2000 09:05:37 GMT"),
	OwnerName:       "Dave Winer",
	OwnerEmail:      "dave@userland.com",
	ExpansionState:  []int{},
	VertScrollState: 1,
	WindowTop:       317,
	WindowLeft:      252,
	WindowBottom:    514,
	WindowRight:     634,
	Outlines: []*Outline{
		{
			Text: "Welcome to Frontier 5!",
			Outlines: []*Outline{
				{Text: "What is Frontier?"},
				{Text: "It's a Content Management System"},
			},
		},
		{
			Text: "Why Manage Content?",
			Outlines: []*Outline{
				{Text: "Form separated from content"},
				{Text: "Make it easy to change the look of a site"},
				{Text: "Keep the technical stuff out of the way of writers"},
				{Text: "Let designers work without having to deal with writers"},
				{Text: "Everyone works on what they do best"},
			},
		},
		{
			Text: "Three groups",
			Outlines: []*Outline{
				{Text: "Writers write"},
				{Text: "Designers design"},
				{Text: "Geeks keep everything working"},
				{Text: "Frontier is for the geeks"},
			},
		},
		{
			Text: "How does Content Flow?",
			Outlines: []*Outline{
				{Text: "Thru LANs, watched, shared folders"},
				{Text: "HTTP Put protocol"},
				{Text: "Email"},
			},
		},
		{
			Text: "Cookie-cutter or workbench?",
			Outlines: []*Outline{
				{Text: "No two organizations work the same way"},
				{Text: "You need a highly customizable environment to make it work"},
				{Text: "A cookie-cutter approach is a dead-end"},
				{Text: "Frontier is <i>designed</i> for customization"},
				{Text: "It's a website system workbench"},
			},
		},
		{
			Text: "Frontier is an environment",
			Outlines: []*Outline{
				{Text: "Everything is integrated"},
				{Text: "Much more powerful"},
				{Text: "Much higher performance"},
				{Text: "Key point!"},
			},
		},
		{
			Text: "The Object Database is the Center",
			Outlines: []*Outline{
				{Text: "Everything is built around a fast scalable object database"},
				{Text: "Millions of hours of burn-in"},
				{Text: "Hierarchical"},
				{Text: "It's also the symbol table for the language"},
			},
		},
		{
			Text: "The scripting language",
			Outlines: []*Outline{
				{Text: "Patterned after C, totally dynamic"},
				{Text: "No need for structure symbols, semicolons or curly braces"},
				{Text: "Because it's integrated with a revolutionary script editor"},
			},
		},
		{
			Text: "The script editor",
			Outlines: []*Outline{
				{Text: "Is an outliner"},
				{Text: "Expand a construct to see the detail"},
				{Text: "Collapse it to hide detail"},
				{Text: "When you move a statement, all the statements under it move too"},
				{Text: "This may be the single most revolutionary feature in Frontier"},
			},
		},
		{
			Text: "Complete script debugger",
			Outlines: []*Outline{
				{Text: "Set a breakpoint"},
				{Text: "Step into and out of procedure calls"},
				{Text: "Easily examine all data while a script is running"},
			},
		},
		{
			Text: "Object oriented website framework",
			Outlines: []*Outline{
				{Text: "Link management with hierarchical glossaries"},
				{Text: "Inherited and overridable attributes"},
				{Text: "Filter scripts also allow overrides and multiple content flows"},
				{Text: "All content is stored in database"},
			},
		},
		{
			Text: "The runtime environment",
			Outlines: []*Outline{
				{Text: "Full built-in TCP support via inetd"},
				{Text: "Fully supports client and server HTTP"},
				{Text: "Fully multi-threaded"},
				{Text: "Large comprehensive verb set"},
				{Text: "Background processes, agents"},
				{Text: "Semaphores"},
			},
		},
		{
			Text: "Editing tools",
			Outlines: []*Outline{
				{Text: "The object database editor is an outliner"},
				{Text: "Outlines are a great format for complex HTML"},
				{Text: "Simple text editor with easy HTML commands"},
			},
		},
		{
			Text: "Key Components of Frontier",
			Outlines: []*Outline{
				{Text: "Integrated database storage system"},
				{Text: "Object oriented website framework"},
				{Text: "Powerful scripting environment with development tools, debugger"},
				{Text: "Outliner and text tools"},
				{Text: "Link management"},
				{Text: "Multithreaded runtime"},
				{Text: "Comprehensive verb set"},
			},
		},
		{
			Text: "Frontier is content management",
			Outlines: []*Outline{
				{Text: "It's not an application development environment"},
				{Text: "It *is* a content management system"},
				{Text: "Suitable for a newspaper or magazine"},
				{Text: "A marketing department"},
				{Text: "A university department"},
			},
		},
		{
			Text: "A brief history of Frontier",
			Outlines: []*Outline{
				{Text: "Automated DTP production with Quark and PageMaker on Mac (1992-93)"},
				{Text: "Transitioned to the web in 1996"},
				{Text: "Ships 1/28/98 for Win32 and Mac"},
				{Text: "Websites are cross platform!"},
				{Text: "So are many utility scripts"},
				{Text: "It's the first truly cross-platform web scripting environment"},
			},
		},
		{
			Text: "A brief future of Frontier",
			Outlines: []*Outline{
				{Text: "Ease of use is our focus"},
				{Text: "Remote procedure calling"},
				{Text: "Sandboxes with scripted firewalls"},
				{Text: "Scalable content"},
				{Text: "XML"},
			},
		},
	},
}

// http://hosting.opml.org/dave/spec/subscriptionList.opml
var subscriptionList = &OPML{
	Version:         "2.0",
	Title:           "mySubscriptions.opml",
	DateCreated:     parseTime("Sat, 18 Jun 2005 12:11:52 GMT"),
	DateModified:    parseTime("Tue, 02 Aug 2005 21:42:48 GMT"),
	OwnerName:       "Dave Winer",
	OwnerEmail:      "dave@scripting.com",
	ExpansionState:  []int{},
	VertScrollState: 1,
	WindowTop:       61,
	WindowLeft:      304,
	WindowBottom:    562,
	WindowRight:     842,
	Outlines: []*Outline{
		{
			Text:        "CNET News.com",
			Description: "Tech news and business reports by CNET News.com. Focused on information technology, core topics include computers, hardware, software, networking, and Internet media.",
			HTMLURL:     parseURL("http://news.com.com/"),
			Language:    "unknown",
			Title:       "CNET News.com",
			Type:        "rss",
			Version:     "RSS2",
			XMLURL:      parseURL("http://news.com.com/2547-1_3-0-5.xml"),
		},
		{
			Text:        "washingtonpost.com - Politics",
			Description: "Politics",
			HTMLURL:     parseURL("http://www.washingtonpost.com/wp-dyn/politics?nav=rss_politics"),
			Language:    "unknown",
			Title:       "washingtonpost.com - Politics",
			Type:        "rss",
			Version:     "RSS2",
			XMLURL:      parseURL("http://www.washingtonpost.com/wp-srv/politics/rssheadlines.xml"),
		},
		{
			Text:        "Scobleizer: Microsoft Geek Blogger",
			Description: "Robert Scoble's look at geek and Microsoft life.",
			HTMLURL:     parseURL("http://radio.weblogs.com/0001011/"),
			Language:    "unknown",
			Title:       "Scobleizer: Microsoft Geek Blogger",
			Type:        "rss",
			Version:     "RSS2",
			XMLURL:      parseURL("http://radio.weblogs.com/0001011/rss.xml"),
		},
		{
			Text:        "Yahoo! News: Technology",
			Description: "Technology",
			HTMLURL:     parseURL("http://news.yahoo.com/news?tmpl=index&cid=738"),
			Language:    "unknown",
			Title:       "Yahoo! News: Technology",
			Type:        "rss",
			Version:     "RSS2",
			XMLURL:      parseURL("http://rss.news.yahoo.com/rss/tech"),
		},
		{
			Text:        "Workbench",
			Description: "Programming and publishing news and comment",
			HTMLURL:     parseURL("http://www.cadenhead.org/workbench/"),
			Language:    "unknown",
			Title:       "Workbench",
			Type:        "rss",
			Version:     "RSS2",
			XMLURL:      parseURL("http://www.cadenhead.org/workbench/rss.xml"),
		},
		{
			Text:        "Christian Science Monitor | Top Stories",
			Description: "Read the front page stories of csmonitor.com.",
			HTMLURL:     parseURL("http://csmonitor.com"),
			Language:    "unknown",
			Title:       "Christian Science Monitor | Top Stories",
			Type:        "rss",
			Version:     "RSS",
			XMLURL:      parseURL("http://www.csmonitor.com/rss/top.rss"),
		},
		{
			Text:        "Dictionary.com Word of the Day",
			Description: "A new word is presented every day with its definition and example sentences from actual published works.",
			HTMLURL:     parseURL("http://dictionary.reference.com/wordoftheday/"),
			Language:    "unknown",
			Title:       "Dictionary.com Word of the Day",
			Type:        "rss",
			Version:     "RSS",
			XMLURL:      parseURL("http://www.dictionary.com/wordoftheday/wotd.rss"),
		},
		{
			Text:        "The Motley Fool",
			Description: "To Educate, Amuse, and Enrich",
			HTMLURL:     parseURL("http://www.fool.com"),
			Language:    "unknown",
			Title:       "The Motley Fool",
			Type:        "rss",
			Version:     "RSS",
			XMLURL:      parseURL("http://www.fool.com/xml/foolnews_rss091.xml"),
		},
		{
			Text:        "InfoWorld: Top News",
			Description: "The latest on Top News from InfoWorld",
			HTMLURL:     parseURL("http://www.infoworld.com/news/index.html"),
			Language:    "unknown",
			Title:       "InfoWorld: Top News",
			Type:        "rss",
			Version:     "RSS2",
			XMLURL:      parseURL("http://www.infoworld.com/rss/news.xml"),
		},
		{
			Text:        "NYT > Business",
			Description: "Find breaking news & business news on Wall Street, media & advertising, international business, banking, interest rates, the stock market, currencies & funds.",
			HTMLURL:     parseURL("http://www.nytimes.com/pages/business/index.html?partner=rssnyt"),
			Language:    "unknown",
			Title:       "NYT > Business",
			Type:        "rss",
			Version:     "RSS2",
			XMLURL:      parseURL("http://www.nytimes.com/services/xml/rss/nyt/Business.xml"),
		},
		{
			Text:        "NYT > Technology",
			Description: "",
			HTMLURL:     parseURL("http://www.nytimes.com/pages/technology/index.html?partner=rssnyt"),
			Language:    "unknown",
			Title:       "NYT > Technology",
			Type:        "rss",
			Version:     "RSS2",
			XMLURL:      parseURL("http://www.nytimes.com/services/xml/rss/nyt/Technology.xml"),
		},
		{
			Text:        "Scripting News",
			Description: "It's even worse than it appears.",
			HTMLURL:     parseURL("http://www.scripting.com/"),
			Language:    "unknown",
			Title:       "Scripting News",
			Type:        "rss",
			Version:     "RSS2",
			XMLURL:      parseURL("http://www.scripting.com/rss.xml"),
		},
		{
			Text:        "Wired News",
			Description: "Technology, and the way we do business, is changing the world we know. Wired News is a technology - and business-oriented news service feeding an intelligent, discerning audience. What role does technology play in the day-to-day living of your life? Wired News tells you. How has evolving technology changed the face of the international business world? Wired News puts you in the picture.",
			HTMLURL:     parseURL("http://www.wired.com/"),
			Language:    "unknown",
			Title:       "Wired News",
			Type:        "rss",
			Version:     "RSS",
			XMLURL:      parseURL("http://www.wired.com/news_drop/netcenter/netcenter.rdf"),
		},
	},
}

// http://hosting.opml.org/dave/spec/states.opml
var states = &OPML{
	Version:         "2.0",
	Title:           "states.opml",
	DateCreated:     parseTime("Tue, 15 Mar 2005 16:35:45 GMT"),
	DateModified:    parseTime("Thu, 14 Jul 2005 23:41:05 GMT"),
	OwnerName:       "Dave Winer",
	OwnerEmail:      "dave@scripting.com",
	ExpansionState:  []int{1, 6, 13, 16, 18, 20},
	VertScrollState: 1,
	WindowTop:       106,
	WindowLeft:      106,
	WindowBottom:    558,
	WindowRight:     479,
	Outlines: []*Outline{
		{
			Text: "United States",
			Outlines: []*Outline{
				{
					Text: "Far West",
					Outlines: []*Outline{
						{Text: "Alaska"},
						{Text: "California"},
						{Text: "Hawaii"},
						{
							Text: "Nevada",
							Outlines: []*Outline{
								{Text: "Reno", Created: parseTime("Tue, 12 Jul 2005 23:56:35 GMT")},
								{Text: "Las Vegas", Created: parseTime("Tue, 12 Jul 2005 23:56:37 GMT")},
								{Text: "Ely", Created: parseTime("Tue, 12 Jul 2005 23:56:39 GMT")},
								{Text: "Gerlach", Created: parseTime("Tue, 12 Jul 2005 23:56:47 GMT")},
							},
						},
						{Text: "Oregon"},
						{Text: "Washington"},
					},
				},
				{
					Text: "Great Plains",
					Outlines: []*Outline{
						{Text: "Kansas"},
						{Text: "Nebraska"},
						{Text: "North Dakota"},
						{Text: "Oklahoma"},
						{Text: "South Dakota"},
					},
				},
				{
					Text: "Mid-Atlantic",
					Outlines: []*Outline{
						{Text: "Delaware"},
						{Text: "Maryland"},
						{Text: "New Jersey"},
						{Text: "New York"},
						{Text: "Pennsylvania"},
					},
				},
				{
					Text: "Midwest",
					Outlines: []*Outline{
						{Text: "Illinois"},
						{Text: "Indiana"},
						{Text: "Iowa"},
						{Text: "Kentucky"},
						{Text: "Michigan"},
						{Text: "Minnesota"},
						{Text: "Missouri"},
						{Text: "Ohio"},
						{Text: "West Virginia"},
						{Text: "Wisconsin"},
					},
				},
				{
					Text: "Mountains",
					Outlines: []*Outline{
						{Text: "Colorado"},
						{Text: "Idaho"},
						{Text: "Montana"},
						{Text: "Utah"},
						{Text: "Wyoming"},
					},
				},
				{
					Text: "New England",
					Outlines: []*Outline{
						{Text: "Connecticut"},
						{Text: "Maine"},
						{Text: "Massachusetts"},
						{Text: "New Hampshire"},
						{Text: "Rhode Island"},
						{Text: "Vermont"},
					},
				},
				{
					Text: "South",
					Outlines: []*Outline{
						{Text: "Alabama"},
						{Text: "Arkansas"},
						{Text: "Florida"},
						{Text: "Georgia"},
						{Text: "Louisiana"},
						{Text: "Mississippi"},
						{Text: "North Carolina"},
						{Text: "South Carolina"},
						{Text: "Tennessee"},
						{Text: "Virginia"},
					},
				},
				{
					Text: "Southwest",
					Outlines: []*Outline{
						{Text: "Arizona"},
						{Text: "New Mexico"},
						{Text: "Texas"},
					},
				},
			},
		},
	},
}

// http://hosting.opml.org/dave/spec/simpleScript.opml
var simpleScript = &OPML{
	Version:         "2.0",
	Title:           "workspace.userlandsamples.doSomeUpstreaming",
	DateCreated:     parseTime("Mon, 11 Feb 2002 22:48:02 GMT"),
	DateModified:    parseTime("Sun, 30 Oct 2005 03:30:17 GMT"),
	OwnerName:       "Dave Winer",
	OwnerEmail:      "dwiner@yahoo.com",
	ExpansionState:  []int{1, 2, 4},
	VertScrollState: 1,
	WindowTop:       74,
	WindowLeft:      41,
	WindowBottom:    314,
	WindowRight:     475,
	Outlines: []*Outline{
		{
			Text:      "Changes",
			IsComment: true,
			Outlines: []*Outline{
				{
					Text: "1/3/02; 4:54:25 PM by DW",
					Outlines: []*Outline{
						{Text: `Change "playlist" to "radio".`},
					},
				},
				{
					Text:      "2/12/01; 1:49:33 PM by DW",
					IsComment: true,
					Outlines: []*Outline{
						{Text: "Test upstreaming by sprinkling a few files in a nice new test folder."},
					},
				},
			},
		},
		{
			Text: "on writetestfile (f, size)",
			Outlines: []*Outline{
				{Text: "file.surefilepath (f)", IsBreakpoint: true},
				{Text: `file.writewholefile (f, string.filledstring ("x", size))`},
			},
		},
		{Text: `local (folder = user.radio.prefs.wwwfolder + "test\\largefiles\\")`},
		{
			Text: "for ch = 'a' to 'z'",
			Outlines: []*Outline{
				{Text: `writetestfile (folder + ch + ".html", random (1000, 16000))`},
			},
		},
	},
}

// http://hosting.opml.org/dave/spec/placesLived.opml
var placesLived = &OPML{
	Version:         "2.0",
	Title:           "placesLived.opml",
	DateCreated:     parseTime("Mon, 27 Feb 2006 12:09:48 GMT"),
	DateModified:    parseTime("Mon, 27 Feb 2006 12:11:44 GMT"),
	OwnerName:       "Dave Winer",
	OwnerID:         parseURL("http://www.opml.org/profiles/sendMail?usernum=1"),
	ExpansionState:  []int{1, 2, 5, 10, 13, 15},
	VertScrollState: 1,
	WindowTop:       242,
	WindowLeft:      329,
	WindowBottom:    665,
	WindowRight:     547,
	Outlines: []*Outline{
		{
			Text: "Places I've lived",
			Outlines: []*Outline{
				{
					Text: "Boston",
					Outlines: []*Outline{
						{Text: "Cambridge"},
						{Text: "West Newton"},
					},
				},
				{
					Text: "Bay Area",
					Outlines: []*Outline{
						{Text: "Mountain View"},
						{Text: "Los Gatos"},
						{Text: "Palo Alto"},
						{Text: "Woodside"},
					},
				},
				{
					Text: "New Orleans",
					Outlines: []*Outline{
						{Text: "Uptown"},
						{Text: "Metairie"},
					},
				},
				{
					Text: "Wisconsin",
					Outlines: []*Outline{
						{Text: "Madison"},
					},
				},
				{
					Text: "Florida",
					Type: "include",
					URL:  parseURL("http://hosting.opml.org/dave/florida.opml"),
				},
				{
					Text: "New York",
					Outlines: []*Outline{
						{Text: "Jackson Heights"},
						{Text: "Flushing"},
						{Text: "The Bronx"},
					},
				},
			},
		},
	},
}

// http://hosting.opml.org/dave/spec/directory.opml
var directory = &OPML{
	Version:         "2.0",
	Title:           "scriptingNewsDirectory.opml",
	DateCreated:     parseTime("Thu, 13 Oct 2005 15:34:07 GMT"),
	DateModified:    parseTime("Tue, 25 Oct 2005 21:33:57 GMT"),
	OwnerName:       "Dave Winer",
	OwnerEmail:      "dwiner@yahoo.com",
	ExpansionState:  []int{},
	VertScrollState: 1,
	WindowTop:       105,
	WindowLeft:      466,
	WindowBottom:    386,
	WindowRight:     964,
	Outlines: []*Outline{
		{
			Text:    "Scripting News sites",
			Created: parseTime("Sun, 16 Oct 2005 05:56:10 GMT"),
			Type:    "link",
			URL:     parseURL("http://hosting.opml.org/dave/mySites.opml"),
		},
		{
			Text:    "News.Com top 100 OPML",
			Created: parseTime("Tue, 25 Oct 2005 21:33:28 GMT"),
			Type:    "link",
			URL:     parseURL("http://news.com.com/html/ne/blogs/CNETNewsBlog100.opml"),
		},
		{
			Text:    "BloggerCon III Blogroll",
			Created: parseTime("Mon, 24 Oct 2005 05:23:52 GMT"),
			Type:    "link",
			URL:     parseURL("http://static.bloggercon.org/iii/blogroll.opml"),
		},
		{
			Text: "TechCrunch reviews",
			Type: "link",
			URL:  parseURL("http://hosting.opml.org/techcrunch.opml.org/TechCrunch.opml"),
		},
		{
			Text: "Tod Maffin's directory of Public Radio podcasts",
			Type: "link",
			URL:  parseURL("http://todmaffin.com/radio.opml"),
		},
		{
			Text: "Adam Curry's iPodder.org directory",
			Type: "link",
			URL:  parseURL("http://homepage.mac.com/dailysourcecode/DSC/ipodderDirectory.opml"),
		},
		{
			Text:    "Memeorandum",
			Created: parseTime("Thu, 13 Oct 2005 15:19:05 GMT"),
			Type:    "link",
			URL:     parseURL("http://tech.memeorandum.com/index.opml"),
		},
		{
			Text:    "DaveNet archive",
			Created: parseTime("Wed, 12 Oct 2005 01:39:56 GMT"),
			Type:    "link",
			URL:     parseURL("http://davenet.opml.org/index.opml"),
		},
	},
}

// http://hosting.opml.org/dave/spec/category.opml
var category = &OPML{
	Version:     "2.0",
	Title:       "Illustrating the category attribute",
	DateCreated: parseTime("Mon, 31 Oct 2005 19:23:00 GMT"),
	Outlines: []*Outline{
		{
			Text: "The Mets are the best team in baseball.",
			Categories: []string{
				"/Philosophy/Baseball/Mets",
				"/Tourism/New York",
			},
			Created: parseTime("Mon, 31 Oct 2005 18:21:33 GMT"),
		},
	},
}

func TestParseSpecification(t *testing.T) {
	testParse(t, "specification.opml", specification)
}

func TestRenderSpecification(t *testing.T) {
	testRender(t, specification)
}

func TestParsePresentation(t *testing.T) {
	testParse(t, "presentation.opml", presentation)
}

func TestRenderPresentation(t *testing.T) {
	testRender(t, presentation)
}

func TestParseSubscriptionList(t *testing.T) {
	testParse(t, "subscriptionList.opml", subscriptionList)
}

func TestRenderSubscriptionList(t *testing.T) {
	testRender(t, subscriptionList)
}

func TestParseStates(t *testing.T) {
	testParse(t, "states.opml", states)
}

func TestRenderStates(t *testing.T) {
	testRender(t, states)
}

func TestParseSimpleScript(t *testing.T) {
	testParse(t, "simpleScript.opml", simpleScript)
}

func TestRenderSimpleScript(t *testing.T) {
	testRender(t, simpleScript)
}

func TestParsePlacesLived(t *testing.T) {
	testParse(t, "placesLived.opml", placesLived)
}

func TestRenderPlacesLived(t *testing.T) {
	testRender(t, placesLived)
}

func TestParseDirectory(t *testing.T) {
	testParse(t, "directory.opml", directory)
}

func TestRenderDirectory(t *testing.T) {
	testRender(t, directory)
}

func TestParseCategory(t *testing.T) {
	testParse(t, "category.opml", category)
}

func TestRenderCategory(t *testing.T) {
	testRender(t, category)
}

func testParse(t *testing.T, filename string, want *OPML) {
	got, err := parseOPML(openTestData(filename))
	if err != nil {
		t.Error("Failed to parse OPML:", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("OPML mismatch\nexpected: %#v\ngot: %#v\n", want, got)
	}
}

func testRender(t *testing.T, want *OPML) {
	var buf bytes.Buffer
	err := Render(&buf, want)
	if err != nil {
		t.Error("Failed to render OPML:", err)
	}

	got, err := Parse(&buf)
	if err != nil {
		t.Error("Failed to parse OPML:", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("OPML mismatch\nexpected: %#v\ngot: %#v\n", want, got)
	}
}

func parseTime(str string) time.Time {
	t, err := time.Parse(time.RFC1123, str)
	if err != nil {
		panic(err)
	}
	return t
}

func parseURL(str string) *url.URL {
	u, err := url.Parse(str)
	if err != nil {
		panic(err)
	}
	return u
}

func openTestData(filename string) io.Reader {
	r, err := os.Open(path.Join("testdata", filename))
	if err != nil {
		panic(err)
	}
	return r
}

func parseOPML(r io.Reader) (*OPML, error) {
	parser := NewParser(r)
	parser.XMLDecoder.CharsetReader = charset.NewReaderLabel
	return parser.Parse()
}

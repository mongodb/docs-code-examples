package tests

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"test-poc/examples/manage-indexes"
	"test-poc/examples/run-queries"
	"testing"
)

func TestAnnQueryBasic(t *testing.T) {
	// Test creating the index and performing a query that relies on the index
	var expected []run_queries.ProjectedMovieResult
	/* Note: we are maintaining different expectations for Atlas and local deployments because not all of the scores match
	 * There are the following discrepancies between scores:
	 * The Time Machine: Atlas: 0.7801066637039185 Local: 0.7801067233085632
	 * Timecop: Atlas: 0.7771612405776978 Local: 0.7771613597869873
	 * Men in Black 3: Atlas: 0.7712380886077881 Local: 0.7712380290031433
	 * Love Story 2050: Atlas: 0.7649372816085815 Local: 0.7649372220039368
	 */
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("no .env file found")
	}
	if os.Getenv("ENV") == "local" {
		expected = []run_queries.ProjectedMovieResult{
			{"Thrill Seekers", "A reporter, learning of time travelers visiting 20th century disasters, tries to change the history they know by averting upcoming disasters.", 0.7892671227455139},
			{"About Time", "At the age of 21, Tim discovers he can travel in time and change what happens and has happened in his own life. His decision to make his world a better place by getting a girlfriend turns out not to be as easy as you might think.", 0.7843604683876038},
			{"The Time Machine", "Hoping to alter the events of the past, a 19th century inventor instead travels 800,000 years into the future, where he finds humankind divided into two warring races.", 0.7801067233085632},
			{"Crusade in Jeans", "After using his mother's newly built time machine, Dolf gets stuck involuntary in the year 1212. He ends up in a children's crusade where he confronts his new friends with modern techniques...", 0.7789170742034912},
			{"Timecop", "An officer for a security agency that regulates time travel, must fend for his life against a shady politician who has a tie to his past.", 0.7771613597869873},
			{"A.P.E.X.", "A time-travel experiment in which a robot probe is sent from the year 2073 to the year 1973 goes terribly wrong thrusting one of the project scientists, a man named Nicholas Sinclair into a...", 0.7730885744094849},
			{"Men in Black 3", "Agent J travels in time to M.I.B.'s early days in 1969 to stop an alien from assassinating his friend Agent K and changing history.", 0.7712380290031433},
			{"Tomorrowland", "Bound by a shared destiny, a teen bursting with scientific curiosity and a former boy-genius inventor embark on a mission to unearth the secrets of a place somewhere in time and space that exists in their collective memory.", 0.7669923901557922},
			{"Love Story 2050", "With the help of his uncle, a man travels to the future to try and bring his girlfriend back to life.", 0.7649372220039368},
			{"The Portal", "A dimension-traveling wizard gets stuck in the 21st century because cell-phone radiation interferes with his magic. With his home world on the brink of war, he seeks help from a jaded ...", 0.7640786170959473},
		}
	} else if os.Getenv("ENV") == "Atlas" {
		expected = []run_queries.ProjectedMovieResult{
			{"Thrill Seekers", "A reporter, learning of time travelers visiting 20th century disasters, tries to change the history they know by averting upcoming disasters.", 0.7892671227455139},
			{"About Time", "At the age of 21, Tim discovers he can travel in time and change what happens and has happened in his own life. His decision to make his world a better place by getting a girlfriend turns out not to be as easy as you might think.", 0.7843604683876038},
			{"The Time Machine", "Hoping to alter the events of the past, a 19th century inventor instead travels 800,000 years into the future, where he finds humankind divided into two warring races.", 0.7801066637039185},
			{"Crusade in Jeans", "After using his mother's newly built time machine, Dolf gets stuck involuntary in the year 1212. He ends up in a children's crusade where he confronts his new friends with modern techniques...", 0.7789170742034912},
			{"Timecop", "An officer for a security agency that regulates time travel, must fend for his life against a shady politician who has a tie to his past.", 0.7771612405776978},
			{"A.P.E.X.", "A time-travel experiment in which a robot probe is sent from the year 2073 to the year 1973 goes terribly wrong thrusting one of the project scientists, a man named Nicholas Sinclair into a...", 0.7730885744094849},
			{"Men in Black 3", "Agent J travels in time to M.I.B.'s early days in 1969 to stop an alien from assassinating his friend Agent K and changing history.", 0.7712380886077881},
			{"Tomorrowland", "Bound by a shared destiny, a teen bursting with scientific curiosity and a former boy-genius inventor embark on a mission to unearth the secrets of a place somewhere in time and space that exists in their collective memory.", 0.7669923901557922},
			{"Love Story 2050", "With the help of his uncle, a man travels to the future to try and bring his girlfriend back to life.", 0.7649372816085815},
			{"The Portal", "A dimension-traveling wizard gets stuck in the 21st century because cell-phone radiation interferes with his magic. With his home world on the brink of war, he seeks help from a jaded ...", 0.7640786170959473},
		}
	} else {
		fmt.Printf("There was no ENV variable set. Ensure your .env file says which environment you're running these tests against.\n")
		t.Fail()
	}

	manage_indexes.ExampleCreateIndexBasic(t)
	var results = run_queries.ExampleAnnBasicQuery(t)

	if VerifyMovieQueryOutput(results, expected) {
		fmt.Printf("The query results match the expected outputs. This test should pass.\n")
	} else {
		t.Fail()
		fmt.Printf("Query results do not match expected query results. This test should fail.\n")
	}
	// Drop the index to clear state for future tests
	manage_indexes.ExampleDropIndex()
}

func TestAnnQueryWithFilter(t *testing.T) {
	var expected []run_queries.ProjectedMovieResultWithFilter
	/* Note: we are maintaining different expectations for Atlas and local deployments because not all of the scores match
	 * There are the following discrepancies between scores:
	 * That Man from Rio: Atlas: 0.7416020035743713 Local: 0.7416019439697266
	 * Willy Wonka & the Chocolate Factory: Atlas: 0.7342107892036438 Local: 0.7342106699943542
	 * Bedknobs and Broomsticks: Atlas: 0.7339356541633606 Local: 0.7339357137680054
	 * Pastoral Hide and Seek: Atlas: 0.733299970626831 Local: 0.7332999110221863
	 * The Three Musketeers: Atlas: 0.7331198453903198 Local: 0.733119785785675
	 */
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("no .env file found")
	}
	if os.Getenv("ENV") == "local" {
		expected = []run_queries.ProjectedMovieResultWithFilter{
			{"Peter Pan", "In this magical tale about the boy who refuses to grow up, Peter Pan and his mischievous fairy sidekick Tinkerbell visit the nursery of Wendy, Michael, and John Darling. With a sprinkling ...", 1960, 0.748110830783844},
			{"Chitty Chitty Bang Bang", "A down-on-his-luck inventor turns a broken-down Grand Prix car into a fancy vehicle for his children, and then they go off on a magical fantasy adventure to save their grandfather in a far-off land.", 1968, 0.7442465424537659},
			{"That Man from Rio", "A young man comes to the rescue of his girlfriend abducted by thieves and brought to Rio. An extravagant adventure ensues.", 1964, 0.7416019439697266},
			{"The Little Prince", "A pilot, stranded in the desert, meets a little boy who is a prince on a planet.", 1974, 0.7378944158554077},
			{"The Red Balloon", "A red balloon with a life of its own follows a little boy around the streets of Paris.", 1956, 0.7342712879180908},
			{"Willy Wonka & the Chocolate Factory", "A poor boy wins the opportunity to tour the most eccentric and wonderful candy factory of all.", 1971, 0.7342106699943542},
			{"Bedknobs and Broomsticks", "An apprentice witch, three kids and a cynical conman search for the missing component to a magic spell useful to the defense of Britain.", 1971, 0.7339357137680054},
			{"Pastoral Hide and Seek", "A young boys' coming of age tale set in a strange, carnivalesque village becomes the recreation of a memory that the director has twenty years later.", 1974, 0.7332999110221863},
			{"The Three Musketeers", "A young swordsman comes to Paris and faces villains, romance, adventure and intrigue with three Musketeer friends.", 1973, 0.733119785785675},
			{"Frosty", "A fairy-tale about a conceited young man and a young woman with a tyrannical step-mother, who must overcome magical trials in order to be together.", 1964, 0.7318308353424072},
		}
	} else if os.Getenv("ENV") == "Atlas" {
		expected = []run_queries.ProjectedMovieResultWithFilter{
			{"Peter Pan", "In this magical tale about the boy who refuses to grow up, Peter Pan and his mischievous fairy sidekick Tinkerbell visit the nursery of Wendy, Michael, and John Darling. With a sprinkling ...", 1960, 0.748110830783844},
			{"Chitty Chitty Bang Bang", "A down-on-his-luck inventor turns a broken-down Grand Prix car into a fancy vehicle for his children, and then they go off on a magical fantasy adventure to save their grandfather in a far-off land.", 1968, 0.7442465424537659},
			{"That Man from Rio", "A young man comes to the rescue of his girlfriend abducted by thieves and brought to Rio. An extravagant adventure ensues.", 1964, 0.7416020035743713},
			{"The Little Prince", "A pilot, stranded in the desert, meets a little boy who is a prince on a planet.", 1974, 0.7378944158554077},
			{"The Red Balloon", "A red balloon with a life of its own follows a little boy around the streets of Paris.", 1956, 0.7342712879180908},
			{"Willy Wonka & the Chocolate Factory", "A poor boy wins the opportunity to tour the most eccentric and wonderful candy factory of all.", 1971, 0.7342107892036438},
			{"Bedknobs and Broomsticks", "An apprentice witch, three kids and a cynical conman search for the missing component to a magic spell useful to the defense of Britain.", 1971, 0.7339356541633606},
			{"Pastoral Hide and Seek", "A young boys' coming of age tale set in a strange, carnivalesque village becomes the recreation of a memory that the director has twenty years later.", 1974, 0.733299970626831},
			{"The Three Musketeers", "A young swordsman comes to Paris and faces villains, romance, adventure and intrigue with three Musketeer friends.", 1973, 0.7331198453903198},
			{"Frosty", "A fairy-tale about a conceited young man and a young woman with a tyrannical step-mother, who must overcome magical trials in order to be together.", 1964, 0.7318308353424072},
		}
	} else {
		fmt.Printf("There was no ENV variable set. Ensure your .env file says which environment you're running these tests against.\n")
		t.FailNow()
	}

	// Test creating the index and performing a query that relies on the index
	manage_indexes.ExampleCreateIndexFilter(t)
	var results = run_queries.ExampleAnnFilterQuery(t)

	if VerifyMovieQueryOutputWithFilter(results, expected) {
		fmt.Printf("The query results match the expected outputs. This test should pass.\n")
	} else {
		t.Fail()
		fmt.Printf("Query results do not match expected query results. This test should fail.\n")
	}

	// Drop the index to clear state for future tests
	manage_indexes.ExampleDropIndex()
}

func TestEnnQuery(t *testing.T) {
	// Test creating the index and performing a query that relies on the index
	var expected []run_queries.ProjectedMovieResult
	/* Note: we are maintaining different expectations for Atlas and local deployments because not all of the scores match
	 * There are the following discrepancies between scores:
	 * When Trumpets Fade: Atlas: 0.7498313188552856 Local: 0.7498312592506409
	 * Saints and Soldiers: Atlas: 0.7435222864151001 Local: 0.7435222268104553
	 * Saints and Soldiers: Atlas: 0.743497371673584 Local: 0.7434973120689392
	 */
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("no .env file found")
	}
	if os.Getenv("ENV") == "local" {
		expected = []run_queries.ProjectedMovieResult{
			{"Red Dawn", "It is the dawn of World War III. In mid-western America, a group of teenagers bands together to defend their town, and their country, from invading Soviet forces.", 0.7700583338737488},
			{"Sands of Iwo Jima", "A dramatization of the World War II Battle of Iwo Jima.", 0.7581185102462769},
			{"White Tiger", "Great Patriotic War, early 1940s. After barely surviving a battle with a mysterious, ghostly-white Tiger tank, Red Army Sergeant Ivan Naydenov becomes obsessed with its destruction.", 0.750884473323822},
			{"P-51 Dragon Fighter", "As World War Two rages on, the allies are about to push the Nazis out of North Africa. That's when the Nazis turn up the heat, unleashing their secret Weapon - dragons.", 0.749922513961792},
			{"When Trumpets Fade", "A private in the latter days of WWII on the German front struggles between his will to survive and what his superiors perceive as a battlefield instinct.", 0.7498312592506409},
			{"Battletruck", "Post World War III futuristic tale of collapsed governments & bankrupt countries heralding a new lawless age.", 0.7497193217277527},
			{"Robot Jox", "It is post-World War III. War is outlawed. In its place, are matches between large Robots called Robot Jox. These matches take place between two large superpowers over disputed territories....", 0.7495121955871582},
			{"The Enemy Below", "During World War II, an American destroyer meets a German U-Boat. Both captains are good ones, and the engagement lasts for a considerable time.", 0.746050238609314},
			{"Saints and Soldiers", "Four American soldiers and one Brit fighting in Europe during World War II struggle to return to Allied territory after being separated from U.S. forces during the historic Malmedy Massacre.", 0.7435222268104553},
			{"Saints and Soldiers", "Four American soldiers and one Brit fighting in Europe during World War II struggle to return to Allied territory after being separated from U.S. forces during the historic Malmedy Massacre.", 0.7434973120689392},
		}
	} else if os.Getenv("ENV") == "Atlas" {
		expected = []run_queries.ProjectedMovieResult{
			{"Red Dawn", "It is the dawn of World War III. In mid-western America, a group of teenagers bands together to defend their town, and their country, from invading Soviet forces.", 0.7700583338737488},
			{"Sands of Iwo Jima", "A dramatization of the World War II Battle of Iwo Jima.", 0.7581185102462769},
			{"White Tiger", "Great Patriotic War, early 1940s. After barely surviving a battle with a mysterious, ghostly-white Tiger tank, Red Army Sergeant Ivan Naydenov becomes obsessed with its destruction.", 0.750884473323822},
			{"P-51 Dragon Fighter", "As World War Two rages on, the allies are about to push the Nazis out of North Africa. That's when the Nazis turn up the heat, unleashing their secret Weapon - dragons.", 0.749922513961792},
			{"When Trumpets Fade", "A private in the latter days of WWII on the German front struggles between his will to survive and what his superiors perceive as a battlefield instinct.", 0.7498313188552856},
			{"Battletruck", "Post World War III futuristic tale of collapsed governments & bankrupt countries heralding a new lawless age.", 0.7497193217277527},
			{"Robot Jox", "It is post-World War III. War is outlawed. In its place, are matches between large Robots called Robot Jox. These matches take place between two large superpowers over disputed territories....", 0.7495121955871582},
			{"The Enemy Below", "During World War II, an American destroyer meets a German U-Boat. Both captains are good ones, and the engagement lasts for a considerable time.", 0.746050238609314},
			{"Saints and Soldiers", "Four American soldiers and one Brit fighting in Europe during World War II struggle to return to Allied territory after being separated from U.S. forces during the historic Malmedy Massacre.", 0.7435222864151001},
			{"Saints and Soldiers", "Four American soldiers and one Brit fighting in Europe during World War II struggle to return to Allied territory after being separated from U.S. forces during the historic Malmedy Massacre.", 0.743497371673584},
		}
	} else {
		fmt.Printf("There was no ENV variable set. Ensure your .env file says which environment you're running these tests against.\n")
		t.FailNow()
	}

	manage_indexes.ExampleCreateIndexBasic(t)
	var results = run_queries.ExampleEnnQuery(t)

	if VerifyMovieQueryOutput(results, expected) {
		fmt.Printf("The query results match the expected outputs. This test should pass.\n")
	} else {
		t.Fail()
		fmt.Printf("Query results do not match expected query results. This test should fail.\n")
	}
	// Drop the index to clear state for future tests
	manage_indexes.ExampleDropIndex()
}

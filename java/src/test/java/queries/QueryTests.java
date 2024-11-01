package queries;

import indexes.CreateIndexBasic;
import indexes.CreateIndexFilter;
import indexes.DropIndex;
import org.bson.Document;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Objects;

import static org.junit.jupiter.api.Assertions.assertEquals;

class QueryTests {
    @AfterEach
    void tearDown() {
        new DropIndex().main(new String[]{"Example placeholder arg"});
    }

    @Test
    @DisplayName("Test basic ANN query")
    void TestAnnQueryBasic() {
        new CreateIndexBasic().main(new String[]{"Example placeholder arg"});
        ArrayList<Document> result = AnnQueryBasic.main(new String[]{"Example placeholder arg"});
        ArrayList<Document> expected = new ArrayList<>();

        /* Note: we are maintaining different expectations for Atlas and local deployments because not all of the scores match
        * There are the following discrepancies between scores:
        * The Time Machine: Atlas: 0.7801066637039185 Local: 0.7801067233085632
        * Timecop: Atlas: 0.7771612405776978 Local: 0.7771613597869873
        * Men in Black 3: Atlas: 0.7712380886077881 Local: 0.7712380290031433
        * Love Story 2050: Atlas: 0.7649372816085815 Local: 0.7649372220039368
        */
        String env = System.getenv("ENV");
        if (Objects.equals(env, "Atlas")) {
            expected.add(new Document("plot", "A reporter, learning of time travelers visiting 20th century disasters, tries to change the history they know by averting upcoming disasters.").append("title", "Thrill Seekers").append("score", 0.7892671227455139));
            expected.add(new Document("plot", "At the age of 21, Tim discovers he can travel in time and change what happens and has happened in his own life. His decision to make his world a better place by getting a girlfriend turns out not to be as easy as you might think.").append("title", "About Time").append("score", 0.7843604683876038));
            expected.add(new Document("plot", "Hoping to alter the events of the past, a 19th century inventor instead travels 800,000 years into the future, where he finds humankind divided into two warring races.").append("title", "The Time Machine").append("score", 0.7801066637039185));
            expected.add(new Document("plot", "After using his mother's newly built time machine, Dolf gets stuck involuntary in the year 1212. He ends up in a children's crusade where he confronts his new friends with modern techniques...").append("title", "Crusade in Jeans").append("score", 0.7789170742034912));
            expected.add(new Document("plot", "An officer for a security agency that regulates time travel, must fend for his life against a shady politician who has a tie to his past.").append("title", "Timecop").append("score", 0.7771612405776978));
            expected.add(new Document("plot", "A time-travel experiment in which a robot probe is sent from the year 2073 to the year 1973 goes terribly wrong thrusting one of the project scientists, a man named Nicholas Sinclair into a...").append("title", "A.P.E.X.").append("score", 0.7730885744094849));
            expected.add(new Document("plot", "Agent J travels in time to M.I.B.'s early days in 1969 to stop an alien from assassinating his friend Agent K and changing history.").append("title", "Men in Black 3").append("score", 0.7712380886077881));
            expected.add(new Document("plot", "Bound by a shared destiny, a teen bursting with scientific curiosity and a former boy-genius inventor embark on a mission to unearth the secrets of a place somewhere in time and space that exists in their collective memory.").append("title", "Tomorrowland").append("score", 0.7669923901557922));
            expected.add(new Document("plot", "With the help of his uncle, a man travels to the future to try and bring his girlfriend back to life.").append("title", "Love Story 2050").append("score", 0.7649372816085815));
            expected.add(new Document("plot", "A dimension-traveling wizard gets stuck in the 21st century because cell-phone radiation interferes with his magic. With his home world on the brink of war, he seeks help from a jaded ...").append("title", "The Portal").append("score", 0.7640786170959473));
        } else if (Objects.equals(env, "local")) {
            expected.add(new Document("plot", "A reporter, learning of time travelers visiting 20th century disasters, tries to change the history they know by averting upcoming disasters.").append("title", "Thrill Seekers").append("score", 0.7892671227455139));
            expected.add(new Document("plot", "At the age of 21, Tim discovers he can travel in time and change what happens and has happened in his own life. His decision to make his world a better place by getting a girlfriend turns out not to be as easy as you might think.").append("title", "About Time").append("score", 0.7843604683876038));
            expected.add(new Document("plot", "Hoping to alter the events of the past, a 19th century inventor instead travels 800,000 years into the future, where he finds humankind divided into two warring races.").append("title", "The Time Machine").append("score", 0.7801067233085632));
            expected.add(new Document("plot", "After using his mother's newly built time machine, Dolf gets stuck involuntary in the year 1212. He ends up in a children's crusade where he confronts his new friends with modern techniques...").append("title", "Crusade in Jeans").append("score", 0.7789170742034912));
            expected.add(new Document("plot", "An officer for a security agency that regulates time travel, must fend for his life against a shady politician who has a tie to his past.").append("title", "Timecop").append("score", 0.7771613597869873));
            expected.add(new Document("plot", "A time-travel experiment in which a robot probe is sent from the year 2073 to the year 1973 goes terribly wrong thrusting one of the project scientists, a man named Nicholas Sinclair into a...").append("title", "A.P.E.X.").append("score", 0.7730885744094849));
            expected.add(new Document("plot", "Agent J travels in time to M.I.B.'s early days in 1969 to stop an alien from assassinating his friend Agent K and changing history.").append("title", "Men in Black 3").append("score", 0.7712380290031433));
            expected.add(new Document("plot", "Bound by a shared destiny, a teen bursting with scientific curiosity and a former boy-genius inventor embark on a mission to unearth the secrets of a place somewhere in time and space that exists in their collective memory.").append("title", "Tomorrowland").append("score", 0.7669923901557922));
            expected.add(new Document("plot", "With the help of his uncle, a man travels to the future to try and bring his girlfriend back to life.").append("title", "Love Story 2050").append("score", 0.7649372220039368));
            expected.add(new Document("plot", "A dimension-traveling wizard gets stuck in the 21st century because cell-phone radiation interferes with his magic. With his home world on the brink of war, he seeks help from a jaded ...").append("title", "The Portal").append("score", 0.7640786170959473));
        }
        assertEquals(expected, result);
    }

    @Test
    @DisplayName("Test ANN query with filter")
    void TestAnnQueryFilter() {
        new CreateIndexFilter().main(new String[]{"Example placeholder arg"});
        ArrayList<Document> result = AnnQueryFilter.main(new String[]{"Example placeholder arg"});
        ArrayList<Document> expected = new ArrayList<>();

        /* Note: we are maintaining different expectations for Atlas and local deployments because not all of the scores match
         * There are the following discrepancies between scores:
         * That Man from Rio: Atlas: 0.7416020035743713 Local: 0.7416019439697266
         * The Red Balloon: Atlas: 0.7342712879180908 Local: 0.734271228313446
         * Willy Wonka & the Chocolate Factory: Atlas: 0.7342107892036438 Local: 0.7342106699943542 CI: 0.734271228313446
         * Bedknobs and Broomsticks: Atlas: 0.7339356541633606 Local: 0.7339357137680054
         * Pastoral Hide and Seek: Atlas: 0.733299970626831 Local: 0.7332999110221863
         * The Three Musketeers: Atlas: 0.7331198453903198 Local: 0.733119785785675
         */
        String env = System.getenv("ENV");
        if (Objects.equals(env, "Atlas")) {
            expected.add(new Document("plot", "In this magical tale about the boy who refuses to grow up, Peter Pan and his mischievous fairy sidekick Tinkerbell visit the nursery of Wendy, Michael, and John Darling. With a sprinkling ...").append("title", "Peter Pan").append("year", 1960).append("score", 0.748110830783844));
            expected.add(new Document("plot", "A down-on-his-luck inventor turns a broken-down Grand Prix car into a fancy vehicle for his children, and then they go off on a magical fantasy adventure to save their grandfather in a far-off land.").append("title", "Chitty Chitty Bang Bang").append("year", 1968).append("score", 0.7442465424537659));
            expected.add(new Document("plot", "A young man comes to the rescue of his girlfriend abducted by thieves and brought to Rio. An extravagant adventure ensues.").append("title", "That Man from Rio").append("year", 1964).append("score", 0.7416020035743713));
            expected.add(new Document("plot", "A pilot, stranded in the desert, meets a little boy who is a prince on a planet.").append("title", "The Little Prince").append("year", 1974).append("score", 0.7378944158554077));
            expected.add(new Document("plot", "A red balloon with a life of its own follows a little boy around the streets of Paris.").append("title", "The Red Balloon").append("year", 1956).append("score", 0.7342712879180908));
            expected.add(new Document("plot", "A poor boy wins the opportunity to tour the most eccentric and wonderful candy factory of all.").append("title", "Willy Wonka & the Chocolate Factory").append("year", 1971).append("score", 0.7342107892036438));
            expected.add(new Document("plot", "An apprentice witch, three kids and a cynical conman search for the missing component to a magic spell useful to the defense of Britain.").append("title", "Bedknobs and Broomsticks").append("year", 1971).append("score", 0.7339356541633606));
            expected.add(new Document("plot", "A young boys' coming of age tale set in a strange, carnivalesque village becomes the recreation of a memory that the director has twenty years later.").append("title", "Pastoral Hide and Seek").append("year", 1974).append("score", 0.733299970626831));
            expected.add(new Document("plot", "A young swordsman comes to Paris and faces villains, romance, adventure and intrigue with three Musketeer friends.").append("title", "The Three Musketeers").append("year", 1973).append("score", 0.7331198453903198));
            expected.add(new Document("plot", "A fairy-tale about a conceited young man and a young woman with a tyrannical step-mother, who must overcome magical trials in order to be together.").append("title", "Frosty").append("year", 1964).append("score", 0.7318308353424072));
        } else if (Objects.equals(env, "local")) {
            expected.add(new Document("plot", "In this magical tale about the boy who refuses to grow up, Peter Pan and his mischievous fairy sidekick Tinkerbell visit the nursery of Wendy, Michael, and John Darling. With a sprinkling ...").append("title", "Peter Pan").append("year", 1960).append("score", 0.748110830783844));
            expected.add(new Document("plot", "A down-on-his-luck inventor turns a broken-down Grand Prix car into a fancy vehicle for his children, and then they go off on a magical fantasy adventure to save their grandfather in a far-off land.").append("title", "Chitty Chitty Bang Bang").append("year", 1968).append("score", 0.7442465424537659));
            expected.add(new Document("plot", "A young man comes to the rescue of his girlfriend abducted by thieves and brought to Rio. An extravagant adventure ensues.").append("title", "That Man from Rio").append("year", 1964).append("score", 0.7416019439697266));
            expected.add(new Document("plot", "A pilot, stranded in the desert, meets a little boy who is a prince on a planet.").append("title", "The Little Prince").append("year", 1974).append("score", 0.7378944158554077));
            expected.add(new Document("plot", "A red balloon with a life of its own follows a little boy around the streets of Paris.").append("title", "The Red Balloon").append("year", 1956).append("score", 0.7342712879180908));
            expected.add(new Document("plot", "A poor boy wins the opportunity to tour the most eccentric and wonderful candy factory of all.").append("title", "Willy Wonka & the Chocolate Factory").append("year", 1971).append("score", 0.7342106699943542));
            expected.add(new Document("plot", "An apprentice witch, three kids and a cynical conman search for the missing component to a magic spell useful to the defense of Britain.").append("title", "Bedknobs and Broomsticks").append("year", 1971).append("score", 0.7339357137680054));
            expected.add(new Document("plot", "A young boys' coming of age tale set in a strange, carnivalesque village becomes the recreation of a memory that the director has twenty years later.").append("title", "Pastoral Hide and Seek").append("year", 1974).append("score", 0.7332999110221863));
            expected.add(new Document("plot", "A young swordsman comes to Paris and faces villains, romance, adventure and intrigue with three Musketeer friends.").append("title", "The Three Musketeers").append("year", 1973).append("score", 0.733119785785675));
            expected.add(new Document("plot", "A fairy-tale about a conceited young man and a young woman with a tyrannical step-mother, who must overcome magical trials in order to be together.").append("title", "Frosty").append("year", 1964).append("score", 0.7318308353424072));
        }
        assertEquals(expected, result);
    }
}
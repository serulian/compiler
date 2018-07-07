$module('match', function () {
  var $static = this;
  this.$class('470dc45c', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.Value = $t.property(function () {
      var $this = this;
      return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Value|3|0e92a8bc": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var firstBool;
    var firstThing;
    var firstValue;
    var fourthBool;
    var fourthThing;
    var fourthValue;
    var secondBool;
    var secondThing;
    var secondValue;
    var thirdBool;
    var thirdThing;
    var thirdValue;
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          firstBool = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
          secondBool = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
          thirdBool = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
          fourthBool = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
          firstValue = $g.match.SomeClass.new();
          secondValue = $t.fastbox(1234, $g.________testlib.basictypes.Integer);
          thirdValue = $t.fastbox('hello world', $g.________testlib.basictypes.String);
          fourthValue = null;
          firstThing = firstValue;
          if ($t.istype(firstThing, $g.match.SomeClass)) {
            $current = 1;
            continue syncloop;
          } else {
            $current = 21;
            continue syncloop;
          }
          break;

        case 1:
          firstBool = firstThing.Value();
          $current = 2;
          continue syncloop;

        case 2:
          secondThing = secondValue;
          if ($t.istype(secondThing, $g.match.SomeClass)) {
            $current = 3;
            continue syncloop;
          } else {
            $current = 17;
            continue syncloop;
          }
          break;

        case 3:
          secondBool = $t.fastbox(!secondThing.Value().$wrapped, $g.________testlib.basictypes.Boolean);
          $current = 4;
          continue syncloop;

        case 4:
          thirdThing = thirdValue;
          if ($t.istype(thirdThing, $g.match.SomeClass)) {
            $current = 5;
            continue syncloop;
          } else {
            $current = 13;
            continue syncloop;
          }
          break;

        case 5:
          thirdBool = $t.fastbox(!thirdThing.Value().$wrapped, $g.________testlib.basictypes.Boolean);
          $current = 6;
          continue syncloop;

        case 6:
          fourthThing = fourthValue;
          if ($t.istype(fourthThing, $g.match.SomeClass)) {
            $current = 7;
            continue syncloop;
          } else {
            $current = 9;
            continue syncloop;
          }
          break;

        case 7:
          fourthBool = $t.fastbox(!fourthThing.Value().$wrapped, $g.________testlib.basictypes.Boolean);
          $current = 8;
          continue syncloop;

        case 8:
          return $t.fastbox(((firstBool.$wrapped && secondBool.$wrapped) && thirdBool.$wrapped) && fourthBool.$wrapped, $g.________testlib.basictypes.Boolean);

        case 9:
          if ($t.istype(fourthThing, $g.________testlib.basictypes.Integer)) {
            $current = 10;
            continue syncloop;
          } else {
            $current = 11;
            continue syncloop;
          }
          break;

        case 10:
          fourthBool = $t.fastbox(fourthThing.$wrapped == 1234, $g.________testlib.basictypes.Boolean);
          $current = 8;
          continue syncloop;

        case 11:
          if (true) {
            $current = 12;
            continue syncloop;
          } else {
            $current = 8;
            continue syncloop;
          }
          break;

        case 12:
          fourthBool = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
          $current = 8;
          continue syncloop;

        case 13:
          if ($t.istype(thirdThing, $g.________testlib.basictypes.Integer)) {
            $current = 14;
            continue syncloop;
          } else {
            $current = 15;
            continue syncloop;
          }
          break;

        case 14:
          thirdBool = $t.fastbox(thirdThing.$wrapped == 1234, $g.________testlib.basictypes.Boolean);
          $current = 6;
          continue syncloop;

        case 15:
          if (true) {
            $current = 16;
            continue syncloop;
          } else {
            $current = 6;
            continue syncloop;
          }
          break;

        case 16:
          thirdBool = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
          $current = 6;
          continue syncloop;

        case 17:
          if ($t.istype(secondThing, $g.________testlib.basictypes.Integer)) {
            $current = 18;
            continue syncloop;
          } else {
            $current = 19;
            continue syncloop;
          }
          break;

        case 18:
          secondBool = $t.fastbox(secondThing.$wrapped == 1234, $g.________testlib.basictypes.Boolean);
          $current = 4;
          continue syncloop;

        case 19:
          if (true) {
            $current = 20;
            continue syncloop;
          } else {
            $current = 4;
            continue syncloop;
          }
          break;

        case 20:
          secondBool = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
          $current = 4;
          continue syncloop;

        case 21:
          if ($t.istype(firstThing, $g.________testlib.basictypes.Integer)) {
            $current = 22;
            continue syncloop;
          } else {
            $current = 23;
            continue syncloop;
          }
          break;

        case 22:
          firstBool = $t.fastbox(firstThing.$wrapped == 4567, $g.________testlib.basictypes.Boolean);
          $current = 2;
          continue syncloop;

        case 23:
          if (true) {
            $current = 24;
            continue syncloop;
          } else {
            $current = 2;
            continue syncloop;
          }
          break;

        case 24:
          firstBool = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
          $current = 2;
          continue syncloop;

        default:
          return;
      }
    }
  };
});

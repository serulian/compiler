$module('match', function () {
  var $static = this;
  this.$class('470dc45c', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $instance.Value = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.fastbox(true, $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Value|3|5ab5941e": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $result;
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
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            firstBool = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
            secondBool = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
            thirdBool = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
            fourthBool = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
            $g.match.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            firstValue = $result;
            secondValue = $t.fastbox(1234, $g.____testlib.basictypes.Integer);
            thirdValue = $t.fastbox('hello world', $g.____testlib.basictypes.String);
            fourthValue = null;
            firstThing = firstValue;
            if ($t.istype(firstThing, $g.match.SomeClass)) {
              $current = 2;
              continue;
            } else {
              $current = 27;
              continue;
            }
            break;

          case 2:
            firstThing.Value().then(function ($result0) {
              firstBool = $result0;
              $result = firstBool;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            $current = 4;
            continue;

          case 4:
            secondThing = secondValue;
            if ($t.istype(secondThing, $g.match.SomeClass)) {
              $current = 5;
              continue;
            } else {
              $current = 23;
              continue;
            }
            break;

          case 5:
            secondThing.Value().then(function ($result0) {
              secondBool = $t.fastbox(!$result0.$wrapped, $g.____testlib.basictypes.Boolean);
              $result = secondBool;
              $current = 6;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 6:
            $current = 7;
            continue;

          case 7:
            thirdThing = thirdValue;
            if ($t.istype(thirdThing, $g.match.SomeClass)) {
              $current = 8;
              continue;
            } else {
              $current = 19;
              continue;
            }
            break;

          case 8:
            thirdThing.Value().then(function ($result0) {
              thirdBool = $t.fastbox(!$result0.$wrapped, $g.____testlib.basictypes.Boolean);
              $result = thirdBool;
              $current = 9;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 9:
            $current = 10;
            continue;

          case 10:
            fourthThing = fourthValue;
            if ($t.istype(fourthThing, $g.match.SomeClass)) {
              $current = 11;
              continue;
            } else {
              $current = 15;
              continue;
            }
            break;

          case 11:
            fourthThing.Value().then(function ($result0) {
              fourthBool = $t.fastbox(!$result0.$wrapped, $g.____testlib.basictypes.Boolean);
              $result = fourthBool;
              $current = 12;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 12:
            $current = 13;
            continue;

          case 13:
            $promise.resolve(firstBool.$wrapped).then(function ($result2) {
              return $promise.resolve($result2 && secondBool.$wrapped).then(function ($result1) {
                return $promise.resolve($result1 && thirdBool.$wrapped).then(function ($result0) {
                  $result = $t.fastbox($result0 && fourthBool.$wrapped, $g.____testlib.basictypes.Boolean);
                  $current = 14;
                  $continue($resolve, $reject);
                  return;
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 14:
            $resolve($result);
            return;

          case 15:
            if ($t.istype(fourthThing, $g.____testlib.basictypes.Integer)) {
              $current = 16;
              continue;
            } else {
              $current = 17;
              continue;
            }
            break;

          case 16:
            fourthBool = $t.fastbox(fourthThing.$wrapped == 1234, $g.____testlib.basictypes.Boolean);
            $current = 13;
            continue;

          case 17:
            if (true) {
              $current = 18;
              continue;
            } else {
              $current = 13;
              continue;
            }
            break;

          case 18:
            fourthBool = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
            $current = 13;
            continue;

          case 19:
            if ($t.istype(thirdThing, $g.____testlib.basictypes.Integer)) {
              $current = 20;
              continue;
            } else {
              $current = 21;
              continue;
            }
            break;

          case 20:
            thirdBool = $t.fastbox(thirdThing.$wrapped == 1234, $g.____testlib.basictypes.Boolean);
            $current = 10;
            continue;

          case 21:
            if (true) {
              $current = 22;
              continue;
            } else {
              $current = 10;
              continue;
            }
            break;

          case 22:
            thirdBool = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
            $current = 10;
            continue;

          case 23:
            if ($t.istype(secondThing, $g.____testlib.basictypes.Integer)) {
              $current = 24;
              continue;
            } else {
              $current = 25;
              continue;
            }
            break;

          case 24:
            secondBool = $t.fastbox(secondThing.$wrapped == 1234, $g.____testlib.basictypes.Boolean);
            $current = 7;
            continue;

          case 25:
            if (true) {
              $current = 26;
              continue;
            } else {
              $current = 7;
              continue;
            }
            break;

          case 26:
            secondBool = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
            $current = 7;
            continue;

          case 27:
            if ($t.istype(firstThing, $g.____testlib.basictypes.Integer)) {
              $current = 28;
              continue;
            } else {
              $current = 29;
              continue;
            }
            break;

          case 28:
            firstBool = $t.fastbox(firstThing.$wrapped == 4567, $g.____testlib.basictypes.Boolean);
            $current = 4;
            continue;

          case 29:
            if (true) {
              $current = 30;
              continue;
            } else {
              $current = 4;
              continue;
            }
            break;

          case 30:
            firstBool = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
            $current = 4;
            continue;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});

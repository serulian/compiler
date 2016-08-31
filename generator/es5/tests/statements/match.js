$module('match', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.Value = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box(true, $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      return $t.createtypesig(['Value', 3, $g.____testlib.basictypes.Boolean.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.match.SomeClass).$typeref()]);
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
            firstBool = $t.box(false, $g.____testlib.basictypes.Boolean);
            secondBool = $t.box(false, $g.____testlib.basictypes.Boolean);
            thirdBool = $t.box(false, $g.____testlib.basictypes.Boolean);
            fourthBool = $t.box(false, $g.____testlib.basictypes.Boolean);
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
            secondValue = $t.box(1234, $g.____testlib.basictypes.Integer);
            thirdValue = $t.box('hello world', $g.____testlib.basictypes.String);
            fourthValue = null;
            firstThing = firstValue;
            if ($t.unbox($t.istype(firstThing, $g.match.SomeClass))) {
              $current = 2;
              continue;
            } else {
              $current = 30;
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
            if ($t.unbox($t.istype(secondThing, $g.match.SomeClass))) {
              $current = 5;
              continue;
            } else {
              $current = 25;
              continue;
            }
            break;

          case 5:
            secondThing.Value().then(function ($result0) {
              secondBool = $t.box(!$t.unbox($result0), $g.____testlib.basictypes.Boolean);
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
            if ($t.unbox($t.istype(thirdThing, $g.match.SomeClass))) {
              $current = 8;
              continue;
            } else {
              $current = 20;
              continue;
            }
            break;

          case 8:
            thirdThing.Value().then(function ($result0) {
              thirdBool = $t.box(!$t.unbox($result0), $g.____testlib.basictypes.Boolean);
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
            if ($t.unbox($t.istype(fourthThing, $g.match.SomeClass))) {
              $current = 11;
              continue;
            } else {
              $current = 15;
              continue;
            }
            break;

          case 11:
            fourthThing.Value().then(function ($result0) {
              fourthBool = $t.box(!$t.unbox($result0), $g.____testlib.basictypes.Boolean);
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
            $promise.resolve($t.unbox(firstBool)).then(function ($result2) {
              return $promise.resolve($result2 && $t.unbox(secondBool)).then(function ($result1) {
                return $promise.resolve($result1 && $t.unbox(thirdBool)).then(function ($result0) {
                  $result = $t.box($result0 && $t.unbox(fourthBool), $g.____testlib.basictypes.Boolean);
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
            if ($t.unbox($t.istype(fourthThing, $g.____testlib.basictypes.Integer))) {
              $current = 16;
              continue;
            } else {
              $current = 18;
              continue;
            }
            break;

          case 16:
            $g.____testlib.basictypes.Integer.$equals(fourthThing, $t.box(1234, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              fourthBool = $result0;
              $result = fourthBool;
              $current = 17;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 17:
            $current = 13;
            continue;

          case 18:
            if (true) {
              $current = 19;
              continue;
            } else {
              $current = 13;
              continue;
            }
            break;

          case 19:
            fourthBool = $t.box(true, $g.____testlib.basictypes.Boolean);
            $current = 13;
            continue;

          case 20:
            if ($t.unbox($t.istype(thirdThing, $g.____testlib.basictypes.Integer))) {
              $current = 21;
              continue;
            } else {
              $current = 23;
              continue;
            }
            break;

          case 21:
            $g.____testlib.basictypes.Integer.$equals(thirdThing, $t.box(1234, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              thirdBool = $result0;
              $result = thirdBool;
              $current = 22;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 22:
            $current = 10;
            continue;

          case 23:
            if (true) {
              $current = 24;
              continue;
            } else {
              $current = 10;
              continue;
            }
            break;

          case 24:
            thirdBool = $t.box(true, $g.____testlib.basictypes.Boolean);
            $current = 10;
            continue;

          case 25:
            if ($t.unbox($t.istype(secondThing, $g.____testlib.basictypes.Integer))) {
              $current = 26;
              continue;
            } else {
              $current = 28;
              continue;
            }
            break;

          case 26:
            $g.____testlib.basictypes.Integer.$equals(secondThing, $t.box(1234, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              secondBool = $result0;
              $result = secondBool;
              $current = 27;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 27:
            $current = 7;
            continue;

          case 28:
            if (true) {
              $current = 29;
              continue;
            } else {
              $current = 7;
              continue;
            }
            break;

          case 29:
            secondBool = $t.box(false, $g.____testlib.basictypes.Boolean);
            $current = 7;
            continue;

          case 30:
            if ($t.unbox($t.istype(firstThing, $g.____testlib.basictypes.Integer))) {
              $current = 31;
              continue;
            } else {
              $current = 33;
              continue;
            }
            break;

          case 31:
            $g.____testlib.basictypes.Integer.$equals(firstThing, $t.box(4567, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              firstBool = $result0;
              $result = firstBool;
              $current = 32;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 32:
            $current = 4;
            continue;

          case 33:
            if (true) {
              $current = 34;
              continue;
            } else {
              $current = 4;
              continue;
            }
            break;

          case 34:
            firstBool = $t.box(false, $g.____testlib.basictypes.Boolean);
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

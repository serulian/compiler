$module('nullcomparecall', function () {
  var $static = this;
  this.$class('9a17167d', 'SomeError', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.Message = $t.property(function () {
      var $this = this;
      return $t.fastbox('huh?', $g.________testlib.basictypes.String);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Message|3|b2b53db7": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.RejectNow = function () {
    throw $g.nullcomparecall.SomeError.new();
  };
  $static.TEST = function () {
    var thing;
    thing = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    return $t.syncnullcompare(thing, function () {
      return $g.nullcomparecall.RejectNow();
    });
  };
});
